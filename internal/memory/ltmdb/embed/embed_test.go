package embed

import (
	"math"
	"testing"

	"github.com/YuruDeveloper/codey/internal/memory/ltmdb/types"
)

const testModelPath = "/home/cecil/LTMDB/public"

func TestNew(t *testing.T) {
	embedder, err := New(testModelPath, 512)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer embedder.Destory()

	if embedder.tokenizer == nil {
		t.Error("tokenizer is nil")
	}
	if embedder.session == nil {
		t.Error("session is nil")
	}
}

func TestEmbed(t *testing.T) {
	embedder, err := New(testModelPath, 512)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer embedder.Destory()

	vector, err := embedder.Embed("hello world")
	if err != nil {
		t.Fatalf("Embed() failed: %v", err)
	}

	// 벡터 길이 확인
	if len(vector) != types.Dimension {
		t.Errorf("vector length: expected %d, got %d", types.Dimension, len(vector))
	}

	// 0이 아닌 값이 있는지 확인
	allZero := true
	for _, v := range vector {
		if v != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("vector is all zeros")
	}
}

func TestEmbedNormalized(t *testing.T) {
	embedder, err := New(testModelPath, 512)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer embedder.Destory()

	vector, err := embedder.Embed("테스트 문장입니다")
	if err != nil {
		t.Fatalf("Embed() failed: %v", err)
	}

	// L2 norm 계산 (정규화되었으면 ≈ 1.0)
	var norm float64
	for _, v := range vector {
		norm += float64(v * v)
	}
	norm = math.Sqrt(norm)

	// 오차 범위 허용 (0.99 ~ 1.01)
	if norm < 0.99 || norm > 1.01 {
		t.Errorf("vector not normalized: L2 norm = %f (expected ~1.0)", norm)
	}
}

func TestEmbedDifferentInputs(t *testing.T) {
	embedder, err := New(testModelPath, 512)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer embedder.Destory()

	vec1, err := embedder.Embed("고양이")
	if err != nil {
		t.Fatalf("Embed() failed: %v", err)
	}

	vec2, err := embedder.Embed("강아지")
	if err != nil {
		t.Fatalf("Embed() failed: %v", err)
	}

	// 다른 입력은 다른 벡터를 생성해야 함
	same := true
	for i := range vec1 {
		if vec1[i] != vec2[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("different inputs produced identical vectors")
	}
}