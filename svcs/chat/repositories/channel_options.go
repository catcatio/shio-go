package repositories

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/catcatio/shio-go/svcs/chat/entities"
)

type ChannelOptionsRepository interface {
	Get(ctx context.Context, channelID string) (*entities.ChannelConfig, error)
}

type channelOptionsRepository struct {
	datastoreClient *datastore.Client
}

func NewChannelOptionsRepository(datastoreClient *datastore.Client) ChannelOptionsRepository {
	return &channelOptionsRepository{
		datastoreClient: datastoreClient,
	}
}

func (c *channelOptionsRepository) Get(ctx context.Context, channelID string) (*entities.ChannelConfig, error) {
	config := new(entities.ChannelConfig)
	key := datastore.NameKey("ChannelConfig", channelID, nil)
	err := c.datastoreClient.Get(ctx, key, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
