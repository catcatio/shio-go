package starter

import (
	"context"
	"github.com/catcatio/shio-go/nub/datastore"
	"github.com/catcatio/shio-go/nub/pubsub"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
)

func LoadLocalOptions() *kernel.ServiceOptions {
	ctx := context.Background()
	log := logger.New("LocalOption")
	log.Debug("enter")
	defer log.Debug("exit")

	gcpProjectID := foundation.EnvStringOrPanic("GCP_PROJECT_ID")
	gcpDataStoreEndpoint := foundation.EnvStringOrPanic("GCP_DATASTORE_ENDPOINT")
	gcpPubsubEndpoint := foundation.EnvStringOrPanic("GCP_PUBSUB_ENDPOINT")

	datastoreClient, err := datastore.NewLocalClient(ctx, gcpProjectID, gcpDataStoreEndpoint)
	if err != nil {
		log.WithError(err).Panic("new datastore client failed")
	}

	pubsubClient, err := pubsub.NewLocalClient(ctx, gcpProjectID, gcpPubsubEndpoint)
	if err != nil {
		log.WithError(err).Panic("new pubsub client failed")
	}

	serviceOptions := &kernel.ServiceOptions{
		PubsubClient:    pubsubClient,
		DatastoreClient: datastoreClient,
	}

	return serviceOptions
}
