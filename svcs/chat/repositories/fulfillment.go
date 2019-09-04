package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	entities2 "github.com/catcatio/shio-go/svcs/chat/entities"
	"net/http"
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

	b, err := json.Marshal(input)
	buf := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", channelConfig.FulfillmentOptions.Endpoint, buf)

	if err != nil {
		panic(err)
	}

	httpClient := http.DefaultClient

	if ctx != nil {
		_, err := httpClient.Do(req.WithContext(ctx))
		if err != nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			default:
			}
		}

		return nil, err
	}
	_, err = httpClient.Do(req)
	return nil, err
}
