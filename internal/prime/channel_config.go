package prime

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/catcatio/shio-go/svcs/chat/entities"
)

func ChannelConfig(datastoreClient *datastore.Client, config *entities.ChannelConfig) (*datastore.Key, error) {
	ctx := context.Background()

	key := datastore.NameKey("ChannelConfig", config.ID, nil)
	outKey, err := datastoreClient.Put(ctx, key, config)

	fmt.Println(err)
	fmt.Println(outKey.String())

	return outKey, err
}
