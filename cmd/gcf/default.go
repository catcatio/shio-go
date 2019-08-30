package gcf

import (
	"context"
	"github.com/catcatio/shio-go/nub/datastore"
	"github.com/catcatio/shio-go/nub/pubsub"
	"github.com/catcatio/shio-go/pkg/kernel"
	pubsub2 "github.com/catcatio/shio-go/pkg/transport/pubsub"
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
		PubsubClients:   pubsub2.NewClients(pubsubClient),
		DatastoreClient: datastoreClient,
	}
}
