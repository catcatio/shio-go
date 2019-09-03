package repositories

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
)

type FulfillmentRepository interface {
	ResolveFulfillment(ctx context.Context, input *entities.IncomingEvent) (*entities.Fulfillment, error)
}

type fulfillmentRepository struct {
	channelConfigRepo ChannelConfigRepository
}

func NewFulfillmentRepository(channelConfigRepo ChannelConfigRepository) FulfillmentRepository {
	return &fulfillmentRepository{channelConfigRepo: channelConfigRepo}
}

func (f *fulfillmentRepository) ResolveFulfillment(ctx context.Context, input *entities.IncomingEvent) (*entities.Fulfillment, error) {
	channelConfig, err := f.channelConfigRepo.Get(ctx, input.ChannelID)
	if err != nil {
		return nil, err
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
