package embed

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/YuruDeveloper/codey/internal/memory/ltmdb/types"
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
	"github.com/viterin/vek/vek32"
	onnx "github.com/yalue/onnxruntime_go"
)

func loadTokenizer(modelPath string, maxLength int) (*tokenizer.Tokenizer, error) {
	tokenizerInstance, err := pretrained.FromFile(filepath.Join(modelPath, "tokenizer.json"))

	if err != nil {
		return nil, err
	}

	configData, err := os.ReadFile(filepath.Join(modelPath, "config.json"))
	if err != nil {
		return nil, err
	}

	var config map[string]any
	err = json.Unmarshal(configData, &config)

	if err != nil {
		return nil, err
	}

	tokenizerConfigData, err := os.ReadFile(filepath.Join(modelPath, "tokenizer_config.json"))
	if err != nil {
		return nil, err
	}

	var tokenizerConfig map[string]any
	err = json.Unmarshal(tokenizerConfigData, &tokenizerConfig)

	if err != nil {
		return nil, err
	}

	tokensMapData, err := os.ReadFile(filepath.Join(modelPath, "special_tokens_map.json"))
	if err != nil {
		return nil, err
	}

	var tokensMap map[string]interface{}
	err = json.Unmarshal(tokensMapData, &tokensMap)

	if err != nil {
		return nil, err
	}

	// Handle overflow when coercing to int, major hassle.
	modelMaxLen := int(min(float64(math.MaxInt32), math.Abs(tokenizerConfig["model_max_length"].(float64))))
	maxLength = min(maxLength, modelMaxLen)

	tokenizerInstance.WithTruncation(&tokenizer.TruncationParams{
		MaxLength: maxLength,
		Strategy:  tokenizer.LongestFirst,
		Stride:    0,
	})

	paddingParams := tokenizer.PaddingParams{
		// Strategy defaults to "BatchLongest"
		Strategy:  *tokenizer.NewPaddingStrategy(),
		Direction: tokenizer.Right,
		PadId:     int(config["pad_token_id"].(float64)),
		PadToken:  tokenizerConfig["pad_token"].(string),
		PadTypeId: 0,
	}
	tokenizerInstance.WithPadding(&paddingParams)

	specialTokens := make([]tokenizer.AddedToken, 0)

	for _, config := range tokensMap {
		switch config := config.(type) {
		case map[string]interface{}:
			{
				specialToken := tokenizer.AddedToken{
					Content:    config["content"].(string),
					SingleWord: config["single_word"].(bool),
					LStrip:     config["lstrip"].(bool),
					RStrip:     config["rstrip"].(bool),
					Normalized: config["normalized"].(bool),
				}
				specialTokens = append(specialTokens, specialToken)
			}
		case string:
			specialToken := tokenizer.AddedToken{
				Content: config,
			}
			specialTokens = append(specialTokens, specialToken)
		default:
			panic(fmt.Sprintf("unknown type for special_tokens_map.json%T", config))
		}
	}
	tokenizerInstance.AddSpecialTokens(specialTokens)

	return tokenizerInstance, nil
}

func encodingToInt32(inputA, inputB, inputC []int) ([]int64, []int64, []int64) {
	if len(inputA) != len(inputB) || len(inputB) != len(inputC) {
		panic("input lengths do not match")
	}
	outputA := make([]int64, len(inputA))
	outputB := make([]int64, len(inputB))
	outputC := make([]int64, len(inputC))
	for i := range inputA {
		outputA[i] = int64(inputA[i])
		outputB[i] = int64(inputB[i])
		outputC[i] = int64(inputC[i])
	}
	return outputA, outputB, outputC
}

func normalize(vector  types.Vector) types.Vector {
	max := float32(0.0)
	for _, value := range vector {
		max += value * value
	}
	max = float32(math.Sqrt(float64(max)))
	min := float32(1e-12)

	normalized := types.Vector(vek32.Div(vector[:],[]float32{(max + min)}))
	
	return normalized
}

type Embed struct {
	tokenizer *tokenizer.Tokenizer
	maxLenght int
	session *onnx.DynamicAdvancedSession
}

func New(path string,maxLenght int) (*Embed , error) {
	onnx.SetSharedLibraryPath(filepath.Join(path,"libonnxruntime.so"))

	if !onnx.IsInitialized() {
		err := onnx.InitializeEnvironment()
		if err != nil {
			return nil , err
		}
	}

	tokenizer , err := loadTokenizer(path,maxLenght)
	if err != nil {
		return nil , err
	}
	session , err := onnx.NewDynamicAdvancedSession(
		filepath.Join(path,"model.onnx"),
		[]string{"input_ids", "attention_mask", "token_type_ids"},
      	[]string{"last_hidden_state"},
	  	nil,
	)
	if err != nil {
		return  nil, err
	}
	return  &Embed{
		tokenizer: tokenizer,
		maxLenght: maxLenght,
		session: session,
	} , nil
}

func (instance *Embed) Destory() error {
	return onnx.DestroyEnvironment()
}

func (instance *Embed) Embed(input string) (types.Vector, error) {
	sequence := tokenizer.NewInputSequence(input)
	encodeInput := tokenizer.NewSingleEncodeInput(sequence)

	encoding, err := instance.tokenizer.Encode(encodeInput, true)
	if err != nil {
		return types.Vector{}, err
	}

	// 토큰 ID 변환
	inputIds, inputMask, inputTypeIds := encodingToInt32(
		encoding.GetIds(),
		encoding.GetAttentionMask(),
		encoding.GetTypeIds(),
	)

	seqLen := int64(len(inputIds))
	inputShape := onnx.NewShape(1, seqLen)

	// 입력 텐서 생성
	inputTensorID, err := onnx.NewTensor(inputShape, inputIds)
	if err != nil {
		return types.Vector{}, err
	}
	defer inputTensorID.Destroy()

	inputTensorMask, err := onnx.NewTensor(inputShape, inputMask)
	if err != nil {
		return types.Vector{}, err
	}
	defer inputTensorMask.Destroy()

	inputTensorType, err := onnx.NewTensor(inputShape, inputTypeIds)
	if err != nil {
		return types.Vector{}, err
	}
	defer inputTensorType.Destroy()

	// 출력 텐서 생성
	outputShape := onnx.NewShape(1, seqLen, int64(types.Dimension))
	outputTensor, err := onnx.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return types.Vector{}, err
	}
	defer outputTensor.Destroy()

	// 세션 실행
	err = instance.session.Run(
		[]onnx.ArbitraryTensor{inputTensorID, inputTensorMask, inputTensorType},
		[]onnx.ArbitraryTensor{outputTensor},
	)
	if err != nil {
		return types.Vector{}, err
	}

	outputData := outputTensor.GetData()
	clsEmbedding := outputData[:types.Dimension]

	return normalize(types.Vector(clsEmbedding)), nil
}