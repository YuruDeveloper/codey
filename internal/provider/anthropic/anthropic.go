package anthropic

import (
	"context"
	"fmt"

	appError "github.com/YuruDeveloper/codey/internal/error"
	"github.com/YuruDeveloper/codey/internal/ports"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/google/uuid"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

var _ ports.Provider = (*Anthropic)(nil)
var _ ports.ClientProvider = (*Anthropic)(nil)

type Anthropic struct {
	client anthropic.Client
	model  int
	datas  []ModelData
}

type ModelData struct {
	Name string
	Id   string
}

const (
    ModelHaikuID   = "claude-haiku-4-5-20251001"
    ModelSonnetID  = "claude-sonnet-4-5-20250929"
    ModelOpusID    = "claude-opus-4-5-20251101"
)

var AnthropicUUID = uuid.New()

var DefaultModels = []ModelData{
    {Name: "haiku4.5", Id: ModelHaikuID},
    {Name: "sonnet4.5", Id: ModelSonnetID},
    {Name: "opus4.5", Id: ModelOpusID},
}

func New(auth ports.Auth) (*Anthropic, error) {
	if dynamic , ok := auth.(ports.DynamicAuth) ; ok {
		if err := dynamic.Update(context.Background()) ; err != nil {
          	return nil, appError.NewError(appError.FailUpdateToken, err)	
		}
	}
	key , authType := auth.Key()
	var client anthropic.Client
	if authType == types.OAuth {
		client = anthropic.NewClient(option.WithAuthToken(key),option.WithHeader("anthropic-beta", "oauth-2025-04-20"))
	} else {
		client = anthropic.NewClient(option.WithAPIKey(key))
	}
	object := &Anthropic{
		client: client,
		model:  0,
	}
	object.getModelsData()
	return object , nil
}

func (instance *Anthropic) Reconnect(auth ports.Auth) error {
	if dynamic , ok := auth.(ports.DynamicAuth) ; ok {
		if err := dynamic.Update(context.Background()) ; err != nil {
      		return appError.NewError(appError.FailUpdateToken, err)
		}	
	}
	key , authType := auth.Key()
	var authOption option.RequestOption
	if authType == types.OAuth {
		instance.client = anthropic.NewClient(option.WithAuthToken(key),option.WithHeader("anthropic-beta", "oauth-2025-04-20"))
	} else {
		instance.client = anthropic.NewClient(option.WithAPIKey(key))
	}
	instance.client = anthropic.NewClient(authOption)
	instance.getModelsData()
	return nil
}

func (instance *Anthropic) getModelsData()  {
	models, err := instance.client.Models.List(context.Background(), anthropic.ModelListParams{})
	if err != nil {
		instance.datas = DefaultModels
		return
	}
	instance.datas = make([]ModelData, len(models.Data))
	for i, model := range models.Data {
		instance.datas[i] = ModelData{
			Name: model.DisplayName,
			Id:   model.ID,
		}
	}
}

func (instance *Anthropic) Models() []string {
	if instance.datas == nil {
		instance.getModelsData()
	}
	names := make([]string, len(instance.datas))
	for i := range len(instance.datas) {
		names[i] = fmt.Sprintf("%s, %s",instance.datas[i].Name,instance.datas[i].Id)  
	}
	return names
}

func (instance *Anthropic) Model() string {
	if instance.datas == nil {
		return ""
	}
	return instance.datas[instance.model].Name
	
}

func (instance *Anthropic) SetModel(index int) {
    if index < 0 || len(instance.datas) <= index {
        instance.model = 0
        return
    }
	instance.model = index
}

func (instance *Anthropic) Send(ctx context.Context, params types.SendParams) (types.Message, error) {
	maxTokens := params.MaxTokens

	if maxTokens == 0 {
		maxTokens = 4096
	}

	message , err := instance.client.Messages.New(ctx,anthropic.MessageNewParams{
		Model: anthropic.Model(instance.datas[instance.model].Id),
		Messages: messageAdapter(params.Messages),
		Tools: toolAdapter(params.Tool),
		MaxTokens: int64(maxTokens),
		System: []anthropic.TextBlockParam{
			{Text: params.SystemPrompt},
		},
	})

	return anthropicToMessageAdapter(message) , err
}

func (instance *Anthropic) GetUUID() uuid.UUID {
	return AnthropicUUID
}
