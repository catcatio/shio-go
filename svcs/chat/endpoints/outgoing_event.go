package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/octofoxio/foundation/logger"
)

type PubsubMessageHandlers map[string]PubsubMessageHandler
type PubsubMessageHandler func(ctx context.Context, m pubsub.RawPubsubMessage) error

func NewOutgoingEventPubsubEndpoint(outgoingEvent usecases.OutgoingEventUsecase, channelConfigRepo repositories.ChannelConfigRepository) PubsubMessageHandler {
	log := logger.New("OutgoingEvent")
	return func(ctx context.Context, m pubsub.RawPubsubMessage) error {

		if m.Data == nil {
			err := fmt.Errorf("data is nil")
			log.WithError(err).Error("handle pubsub failed")
			return err
		}

		input := new(entities.OutgoingEvent)

		if err := json.Unmarshal(m.Data, input); err != nil {
			log.WithError(err).WithField("data", string(m.Data)).Error("unmarshal data failed")
			return err
		}

		ctx = shio.AppendRequestIDToContext(ctx, input.RequestID)
		channelConfig, err := channelConfigRepo.Get(ctx, input.ChannelID)
		if err != nil {
			log.WithError(err).Error("get channel config failed")
			return err
		}

		ctx = shio.PutContextValue(ctx, "channel_config", *channelConfig)

		log.Println(input)

		return outgoingEvent.Handle(ctx, input)
	}
}
