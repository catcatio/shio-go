package usecases

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/nub"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	lu "github.com/catcatio/shio-go/utils/line"
)

type Line interface {
	RequestParser() lu.RequestParser
	HandleIncomingEvents(ctx context.Context, in *WebhookInput) error
}

type line struct {
	channelSecret      string
	channelAccessToken string
	parser             lu.RequestParser
	intent             repositories.IntentRepository
}

func New(channelSecret string, channelAccessToken string, parser lu.RequestParser, intent repositories.IntentRepository) Line {
	return &line{channelSecret: channelSecret, channelAccessToken: channelAccessToken, parser: parser}
}

func (l *line) RequestParser() lu.RequestParser {
	return l.parser

}

func (l *line) HandleIncomingEvents(ctx context.Context, in *WebhookInput) (err error) {
	fmt.Println(in.Provider)
	fmt.Println(in.Events)

	parsedEvents := make([]*entities.ParsedEvent, 0)

	for _, e := range in.Events {
		//intent, err := l.intent.Detect(e.GetMessage())
		//
		//if err != nil {
		//	intent = nil
		//	fmt.Println(err)
		//}

		parsedEvent := &entities.ParsedEvent{
			RequestID:    nub.NewID(),
			Message:      e.GetMessage(),
			ProviderType: e.GetProvider(),
			ReplyToken:   e.GetReplyToken(),
			TimeStamp:    e.GetTimestamp(),
			Source:       e.GetSource(),
			Original:     e.GetOriginalEvent(),
			//Intent:       intent,
		}

		parsedEvents = append(parsedEvents, parsedEvent)
	}

	fmt.Println(parsedEvents)
	return nil
}

type WebhookInput struct {
	Provider entities.ProviderType
	Events   []entities.Event `json:"events"`
}

type WebhookOutput struct {
	ParsedEvents []*entities.ParsedEvent
}
