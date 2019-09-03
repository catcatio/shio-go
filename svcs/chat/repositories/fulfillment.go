package repositories

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	entities2 "github.com/catcatio/shio-go/svcs/chat/entities"
)

type FulfillmentRepository interface {
	ResolveFulfillment(ctx context.Context, input *entities.IncomingEvent) (*entities.Fulfillment, error)
}

type fulfillmentRepository struct {
}

func NewFulfillmentRepository() FulfillmentRepository {
	return &fulfillmentRepository{}
}

func (f *fulfillmentRepository) ResolveFulfillment(ctx context.Context, input *entities.IncomingEvent) (*entities.Fulfillment, error) {
	var channelConfig entities2.ChannelConfig
	if config, ok := ctx.Value("channel_config").(entities2.ChannelConfig); !ok {
		return nil, ErrChannelConfigNotSet
	} else {
		channelConfig = config
	}

	fmt.Println(input.Intent)

	fmt.Println(channelConfig.FulfillmentOptions.Endpoint)

	return &entities.Fulfillment{
		Name: "dummy",
		Parameters: entities.Parameters{
			"yes": "no",
		},
	}, nil
}
