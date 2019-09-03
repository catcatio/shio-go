package endpoints

import (
	"context"
	"encoding/json"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/octofoxio/foundation/logger"
)

func NewFulfillmentPubsubEndpoint(fulfillment usecases.FulfillmentUsecase) PubsubMessageHandler {
	log := logger.New("NewFulfillmentPubsubEndpoint").WithServiceInfo("handle")
	return func(ctx context.Context, m pubsub.RawPubsubMessage) error {
		input := new(entities.IncomingEvent)

		if err := json.Unmarshal(m.Data, input); err != nil {
			log.WithError(err).WithField("_data", string(m.Data)).Error("unmarshal data failed")
			return err
		}

		ctx = shio.AppendRequestIDToContext(ctx, input.RequestID)

		return fulfillment.Dispatch(ctx, input)
	}
}