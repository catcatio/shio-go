package usecases

import (
	"context"
	"github.com/catcatio/shio-go/nub"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	lu "github.com/catcatio/shio-go/utils/line"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
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

	log *logger.Logger
}

func New(lineOptions *kernel.LineChatOptions, parser lu.RequestParser, intent repositories.IntentRepository) Line {
	return &line{
		channelSecret:      lineOptions.ChannelSecret,
		channelAccessToken: lineOptions.ChannelAccessToken,
		parser:             parser,
		intent:             intent,
		log:                logger.New("LineUsecase"),
	}
}

func (l *line) RequestParser() lu.RequestParser {
	return l.parser

}

func (l *line) HandleIncomingEvents(ctx context.Context, in *WebhookInput) (err error) {
	log := l.log.WithServiceInfo("HandleIncomingEvents").WithRequestID(foundation.GetRequestIDFromContext(ctx))
	log.Infof("%d incoming event(s) from %s", len(in.Events), in.Provider)

	parsedEvents := make([]*entities.ParsedEvent, 0)

	for _, e := range in.Events {
		intent, err := l.intent.Detect(ctx, e.GetMessage(), "en")

		if err != nil {
			intent = nil
			log.Println(err)
		}

		parsedEvent := &entities.ParsedEvent{
			RequestID:    nub.NewID(),
			Message:      e.GetMessage(),
			ProviderType: e.GetProvider(),
			ReplyToken:   e.GetReplyToken(),
			TimeStamp:    e.GetTimestamp(),
			Source:       e.GetSource(),
			Original:     e.GetOriginalEvent(),
			Intent:       intent,
		}

		parsedEvents = append(parsedEvents, parsedEvent)
	}

	log.Println(parsedEvents)
	return nil
}

type WebhookInput struct {
	Provider entities.ProviderType
	Events   []entities.Event `json:"events"`
}

type WebhookOutput struct {
	ParsedEvents []*entities.ParsedEvent
}
