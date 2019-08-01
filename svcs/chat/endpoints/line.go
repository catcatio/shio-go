package endpoints

import (
	"fmt"
	"github.com/catcatio/shio-go/nub"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
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

func NewLineWebhookEndpointFunc(line usecases.Line) EndpointFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := nub.NewID()
		log := logger.New("endpoint").WithServiceInfo("webhook").WithRequestID(reqID)
		events, _, err := line.RequestParser().Parse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "err: %s", err.Error())
			log.WithError(err).Error("line.RequestParser failed")
			return
		}

		if len(events) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Sprintf("err: %s", http.StatusText(http.StatusBadRequest))
			_, _ = fmt.Fprintf(w, msg)
			log.Error(msg)
			return
		}

		ctx := shio.NewContextWithRequestID(reqID)
		le := make([]entities.Event, 0)
		for _, e := range events {
			if isSystemMessage(e) {
				log.Infof("system message received")
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, "%s", "OK")
				return // response and return
			}

			le = append(le, &entities.LineEvent{Event: e})
		}

		err = line.HandleIncomingEvents(ctx, &usecases.WebhookInput{
			Provider: entities.ProviderTypeLine,
			Events:   le,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "err: %s", err.Error())
			log.WithError(err).Error("handle incoming events failed")
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "OK")
	}
}
