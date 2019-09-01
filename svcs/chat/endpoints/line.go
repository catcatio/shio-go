package endpoints

import (
	"context"
	"github.com/catcatio/shio-go/nub"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	chatentities "github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	lineutils "github.com/catcatio/shio-go/utils/line"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/octofoxio/foundation/logger"
	"net/http"
)

func isSystemMessage(event *linebot.Event) bool {
	if event == nil {
		return false
	}

	if event.Type != linebot.EventTypeMessage {
		return false
	}

	return event.ReplyToken == "00000000000000000000000000000000" || event.ReplyToken == "ffffffffffffffffffffffffffffffff"
}

func newLineEndpointFunc(incomingEvent usecases.IncomingEventUsecase) ProviderEndpointFunc {
	return func(ctx context.Context, options *chatentities.ChannelConfig, w http.ResponseWriter, r *http.Request) {
		log := logger.New("newLineEndpointFunc")
		lineOptions := options.LineChatOptions
		requestParser := lineutils.NewParser(lineOptions.ChannelSecret)

		// parse events and validate signature
		events, _, err := requestParser.Parse(r)
		if err == lineutils.ErrInvalidSignature {
			writeBadRequest(w)
			log.WithError(err).Error("validation failed")
			return
		} else if err != nil {
			writeErrorResponse(w, err.Error())
			log.WithError(err).Error("parse events failed")
			return
		}

		if len(events) <= 0 {
			log.Error("empty event, just stop")
			writeBadRequest(w)
			return
		}

		le := make([]*entities.IncomingEvent, 0)
		for _, e := range events {
			if isSystemMessage(e) {
				log.Infof("system message received")
				writeOKResponse(w)
				return // response and return
			}

			le = append(le, (&entities.LineEvent{Event: e}).IncomingEvent(nub.NewID(), options.ID))
		}

		// just return 200 OK, handle the rest internally
		writeOKResponse(w)

		// pass events to handler
		go func(ctx context.Context, log *logger.Logger) {
			defer func() {
				var r interface{}
				if r = recover(); r == nil {
					return
				}

				if e, ok := r.(error); ok {
					log.WithError(e).Error("panic occurred")
				}
			}()
			err = incomingEvent.HandleEvents(ctx,
				&chatentities.WebhookInput{
					Provider:  entities.ProviderTypeLine,
					ChannelID: options.ID,
					Events:    le,
				})
		}(ctx, log)
	}
}
