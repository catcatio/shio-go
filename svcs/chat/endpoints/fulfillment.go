package endpoints

import (
	"context"
	"encoding/json"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/octofoxio/foundation/logger"
)

func NewFulfillmentPubsubEndpoint(fulfillment usecases.FulfillmentUsecase, channelConfigRepo repositories.ChannelConfigRepository) PubsubMessageHandler {
	log := logger.New("NewFulfillmentPubsubEndpoint").WithServiceInfo("handle")
	return func(ctx context.Context, m pubsub.RawPubsubMessage) error {
		input := new(entities.IncomingEvent)

		if err := json.Unmarshal(m.Data, input); err != nil {
			log.WithError(err).WithField("_data", string(m.Data)).Error("unmarshal data failed")
			return err
		}

		ctx = shio.AppendRequestIDToContext(ctx, input.RequestID)
		channelConfig, err := channelConfigRepo.Get(ctx, input.ChannelID)
		if err != nil {
			log.WithError(err).Error("get channel config failed")
			return err
		}

		ctx = shio.PutContextValue(ctx, "channel_config", *channelConfig)

		return fulfillment.Dispatch(ctx, input)
	}
}
