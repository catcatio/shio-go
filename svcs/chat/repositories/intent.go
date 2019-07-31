package repositories

import (
	"cloud.google.com/go/dialogflow/apiv2"
	"context"
	"errors"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var (
	ErrMessageTypeNotSupported = errors.New("message type not supported")
)

type IntentRepository interface {
	Detect(ctx context.Context, m entities.Message, languageCode string) (*entities.Intent, error)
}

type dialogflowIntentRepository struct {
	projectName string
	client      *dialogflow.SessionsClient
}

func NewDialogflowIntentRepository(projectName string, credentialsJson string) IntentRepository {
	ctx := context.Background()
	opts := option.WithCredentialsJSON([]byte(credentialsJson))
	client, err := dialogflow.NewSessionsClient(ctx, opts)

	if err != nil {
		panic(err)
	}

	return &dialogflowIntentRepository{
		projectName: projectName,
		client:      client,
	}
}

func (i *dialogflowIntentRepository) Detect(ctx context.Context, m entities.Message, languageCode string) (*entities.Intent, error) {
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

	response, err := i.client.DetectIntent(ctx, &dialogflowpb.DetectIntentRequest{
		Session: "session",
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

	return &entities.Intent{
		Name: response.QueryResult.GetAction(),
		//Text: response.QueryResult.Parameters
	}, nil
}
