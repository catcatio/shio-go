package repositories

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/catcatio/shio-go/svcs/chat/entities"
)

type ChannelConfigRepository interface {
	Get(ctx context.Context, channelID string) (*entities.ChannelConfig, error)
}

type channelConfigRepository struct {
	datastoreClient *datastore.Client
}

func NewChannelConfigRepository(datastoreClient *datastore.Client) ChannelConfigRepository {
	return &channelConfigRepository{
		datastoreClient: datastoreClient,
	}
}

func (c *channelConfigRepository) Get(ctx context.Context, channelID string) (*entities.ChannelConfig, error) {
	config := new(entities.ChannelConfig)
	key := datastore.NameKey("ChannelConfig", channelID, nil)
	err := c.datastoreClient.Get(ctx, key, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
