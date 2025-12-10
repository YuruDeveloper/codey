package anthropic

import (
	"context"
	"fmt"

	"github.com/YuruDeveloper/codey/internal/auth"
	appError "github.com/YuruDeveloper/codey/internal/error"
	"github.com/YuruDeveloper/codey/internal/provider"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

var _ provider.Provider = (*Anthropic)(nil)
var _ provider.ClientProvider = (*Anthropic)(nil)

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

var DefaultModels = []ModelData{
    {Name: "haiku4.5", Id: ModelHaikuID},
    {Name: "sonnet4.5", Id: ModelSonnetID},
    {Name: "opus4.5", Id: ModelOpusID},
}

func New(key auth.Auth) (*Anthropic, error) {
	if dynamic , ok := key.(auth.DynamicAuth) ; ok {
		if err := dynamic.Update(context.Background()) ; err != nil {
          	return nil, appError.NewError(appError.FailUpdateToken, err)	
		}
	}
	client := anthropic.NewClient(option.WithAPIKey(key.Key()))
	object := &Anthropic{
		client: client,
		model:  0,
	}
	object.getModelsData()
	return object , nil
}

func (instance *Anthropic) Reconnect(key auth.Auth) error {
	if dynamic , ok := key.(auth.DynamicAuth) ; ok {
		if err := dynamic.Update(context.Background()) ; err != nil {
      		return appError.NewError(appError.FailUpdateToken, err)
		}	
	}
	instance.client = anthropic.NewClient(option.WithAPIKey(key.Key()))	
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
