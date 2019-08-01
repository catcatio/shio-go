package repositories

import (
	"cloud.google.com/go/dialogflow/apiv2"
	"context"
	"errors"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/octofoxio/foundation/logger"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var (
	ErrMessageTypeNotSupported = errors.New("message type not supported")
)

type IntentRepository interface {
	Detect(ctx context.Context, m entities.Message, sessionID string, languageCode string) (*entities.Intent, error)
}

type dialogflowIntentRepository struct {
	projectName string
	client      *dialogflow.SessionsClient
	log         *logger.Logger
}

func NewDialogflowIntentRepository(projectName string, credentialsJson string) IntentRepository {
	log := logger.New("dialogflowIntent")
	ctx := context.Background()
	opts := option.WithCredentialsJSON([]byte(credentialsJson))
	client, err := dialogflow.NewSessionsClient(ctx, opts)

	if err != nil {
		// panic here
		log.WithError(err).WithServiceInfo("New").Panic("new dialog client failed")
	}

	return &dialogflowIntentRepository{
		projectName: projectName,
		client:      client,
		log:         log,
	}
}

func (i *dialogflowIntentRepository) Detect(ctx context.Context, m entities.Message, sessionID string, languageCode string) (*entities.Intent, error) {
	if m.GetType() != entities.MessageTypeTextMessage {
		return nil, ErrMessageTypeNotSupported
	}

	var eventMsg *entities.EventMessage
	if e, ok := m.(*entities.EventMessage); !ok {
		return nil, ErrMessageTypeNotSupported
	} else {
		eventMsg = e
	}

	textMsg := eventMsg.Parameters.String("Text")
	session := fmt.Sprintf("projects/%s/agent/sessions/%s", i.projectName, sessionID)
	response, err := i.client.DetectIntent(ctx, &dialogflowpb.DetectIntentRequest{
		Session: session,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
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
