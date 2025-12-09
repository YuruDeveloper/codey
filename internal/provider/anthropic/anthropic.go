package anthropic

import (
	"context"
	"fmt"

	"github.com/YuruDeveloper/codey/internal/auth"
	"github.com/YuruDeveloper/codey/internal/provider"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

var _ provider.Provider = (*Anthropic)(nil)

type Anthropic struct {
	client anthropic.Client
	model  int
	datas  []ModelData
}

type ModelData struct {
	Name string
	Id   string
}

func New(auth auth.Auth) *Anthropic {
	auth.Update(context.Background())
	client := anthropic.NewClient(option.WithAPIKey(auth.Key()))
	return &Anthropic{
		client: client,
		model:  0,
	}
}

func (instance *Anthropic) getModelsData()  {
	models, err := instance.client.Models.List(context.Background(), anthropic.ModelListParams{})
	if err != nil {
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
	instance.model = index
}

func (instance *Anthropic) Send(ctx context.Context, params provider.SendParams) <-chan types.StreamEvent {
	eventChan := make(chan types.StreamEvent)
	stream(
		ctx,
		&instance.client,
		anthropic.Model(instance.datas[instance.model].Id),
		eventChan,
		initStream(params),
	)
	return eventChan
}
