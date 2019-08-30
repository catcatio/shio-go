package endpoints

import (
	"context"
	"encoding/json"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/octofoxio/foundation/logger"
)

func NewIncomingEventPubsubEndpoint(intent usecases.IntentUsecase) PubsubMessageHandler {
	log := logger.New("IncomingEvent")
	return func(ctx context.Context, m pubsub.RawPubsubMessage) error {
		input := new(entities.IncomingEvent)

		if err := json.Unmarshal(m.Data, input); err != nil {
			log.WithError(err).WithField("data", string(m.Data)).Error("unmarshal data failed")
			return err
		}

		return intent.HandleEvents(ctx, input)
	}
}
