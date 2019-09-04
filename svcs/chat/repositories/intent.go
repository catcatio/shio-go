package repositories

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"errors"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	entities2 "github.com/catcatio/shio-go/svcs/chat/entities"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/option"
	dialogflow2 "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"strings"
)

var (
	ErrMessageTypeNotSupported = errors.New("message type not supported")
	ErrChannelConfigNotSet     = errors.New("channel config not set")
	ErrInvalidConfig           = errors.New("invalid configuration")
)

type IntentRepository interface {
	AddDetector(detector IntentDetector)
	Detect(ctx context.Context, event *entities.IncomingEvent) (*entities.Intent, error)
}

type intentRepository struct {
	pubsubChannelRepository PubsubChannelRepository
	intentDetectorProvider  IntentDetectorProvider
}

func NewIntentRepository() IntentRepository {
	intent := &intentRepository{
		intentDetectorProvider: newIntentDetectorProvider(),
	}

	// to add other provider here
	intent.AddDetector(NewDialogflowIntentDetector())
	return intent
}

func (i *intentRepository) AddDetector(detector IntentDetector) {
	i.intentDetectorProvider.Add(detector)
}

func (i *intentRepository) Detect(ctx context.Context, event *entities.IncomingEvent) (*entities.Intent, error) {
	var channelConfig entities2.ChannelConfig
	if config, ok := ctx.Value("channel_config").(entities2.ChannelConfig); !ok {
		return nil, ErrChannelConfigNotSet
	} else {
		channelConfig = config
	}

	detector, err := i.intentDetectorProvider.Get(channelConfig.IntentDetector)

	if err != nil {
		return nil, err
	}
	return detector.Detect(ctx, &channelConfig, event)
}

type IntentDetectorProvider interface {
	Add(detector IntentDetector)
	Get(provider string) (IntentDetector, error)
}

type IntentDetector interface {
	Name() string
	Detect(ctx context.Context, channelConfig *entities2.ChannelConfig, event *entities.IncomingEvent) (*entities.Intent, error)
}

type defaultIntentDetectorProvider struct {
	detectors map[string]IntentDetector
}

func (d *defaultIntentDetectorProvider) Add(detector IntentDetector) {
	d.detectors[detector.Name()] = detector
}

func (d *defaultIntentDetectorProvider) Get(name string) (IntentDetector, error) {
	if d, ok := d.detectors[name]; ok {
		return d, nil
	}

	return nil, errors.New("not found")
}

func newIntentDetectorProvider() IntentDetectorProvider {
	return &defaultIntentDetectorProvider{
		detectors: make(map[string]IntentDetector),
	}
}

type dialogflowIntentDetector struct {
}

func NewDialogflowIntentDetector() IntentDetector {
	return &dialogflowIntentDetector{}
}

func (d *dialogflowIntentDetector) Name() string {
	return "dialogflow"
}

func (d *dialogflowIntentDetector) Detect(ctx context.Context, channelConfig *entities2.ChannelConfig, event *entities.IncomingEvent) (*entities.Intent, error) {
	eventMsg := event.Message
	if eventMsg.GetType() != entities.MessageTypeTextMessage {
		return nil, ErrMessageTypeNotSupported
	}

	if channelConfig.DialogflowOptions == nil ||
		channelConfig.DialogflowOptions.CredentialsJson == "" ||
		channelConfig.DialogflowOptions.ProjectID == "" {
		return nil, ErrInvalidConfig
	}

	dfConfig := channelConfig.DialogflowOptions

	credentialsJson := strings.Trim(dfConfig.CredentialsJson, `"`)
	credentialsJson = strings.Replace(credentialsJson, `'`, `"`, -1) // use string.Replace for the sake of go 1.11 in gcf

	opts := option.WithCredentialsJSON([]byte(credentialsJson))
	client, err := dialogflow.NewSessionsClient(ctx, opts)

	if err != nil {
		return nil, err
	}

	// TODO: detect language (set as option)
	languageCode := "en"
	sessionID := event.Source.GetSessionID()
	textMsg := eventMsg.Parameters.String("Text")
	session := fmt.Sprintf("projects/%s/agent/sessions/%s", dfConfig.ProjectID, sessionID)
	response, err := client.DetectIntent(ctx, &dialogflow2.DetectIntentRequest{
		Session: session,
		QueryInput: &dialogflow2.QueryInput{
			Input: &dialogflow2.QueryInput_Text{
				Text: &dialogflow2.TextInput{
					Text:         textMsg,
					LanguageCode: languageCode,
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	params := make(entities.Parameters)
	for k, v := range response.QueryResult.Parameters.Fields {
		switch x := v.Kind.(type) {
		case *structpb.Value_NullValue:
			params[k] = x.NullValue
			break
		case *structpb.Value_NumberValue:
			params[k] = x.NumberValue
			break
		case *structpb.Value_StringValue:
			params[k] = x.StringValue
			break
		case *structpb.Value_BoolValue:
			params[k] = x.BoolValue
			break
		case *structpb.Value_StructValue:
			params[k] = x.StructValue
			break
		case *structpb.Value_ListValue:
			params[k] = x.ListValue
			break

		}
	}
	return &entities.Intent{
		Name:                     response.QueryResult.GetAction(),
		FulfillmentText:          response.QueryResult.GetFulfillmentText(),
		Parameters:               params,
		AllRequiredParamsPresent: response.QueryResult.GetAllRequiredParamsPresent(),
	}, nil
}
