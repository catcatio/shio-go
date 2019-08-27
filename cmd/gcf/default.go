package gcf

import (
	"context"
	"github.com/catcatio/shio-go/nub/datastore"
	"github.com/catcatio/shio-go/nub/pubsub"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/octofoxio/foundation/logger"
)

func newServiceOptions(ctx context.Context, projectID string) *kernel.ServiceOptions {
	log := logger.New("gcf").WithServiceInfo("newServiceOptions")
	pubsubClient, err := pubsub.NewDefaultClient(ctx, projectID)
	if err != nil {
		log.Panicf("create pubsub client failed %s", err.Error())
	}

	datastoreClient, err := datastore.NewDefaultClient(ctx, projectID)
	if err != nil {
		log.Panicf("create datastore client failed %s", err.Error())
	}

	return &kernel.ServiceOptions{
		PubsubClient:    pubsubClient,
		DatastoreClient: datastoreClient,
	}
}
