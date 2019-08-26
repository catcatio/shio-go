package test

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/cmd/shio-svc/starter"
	"github.com/catcatio/shio-go/nub/datastore"
	"github.com/catcatio/shio-go/nub/pubsub"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/octofoxio/foundation/logger"
	"os"
	"testing"
	"time"
)

type IntegrationTestOption struct {
	Name        string
	Host        string
	Port        string
	StartServer bool
}

func IntegrationTestSetup(m *testing.M, options *IntegrationTestOption) {
	if testing.Short() {
		fmt.Println("skipping setup integration test")
		m.Run()
	}

	var service *starter.Service
	defer func() {
		time.Sleep(time.Second)
		c := m.Run()

		time.Sleep(time.Second)
		if service != nil {
			service.Stop()
		}

		os.Exit(c)
	}()

	if options.StartServer {
		serviceOptions := makeServiceOptions(options)
		s := starter.Starter{
			Host:    options.Host,
			Port:    options.Port,
			Options: serviceOptions,
			Log:     logger.New(fmt.Sprintf("%s_test", options.Name)),
		}

		service = s.Start()
	}
}

func makeServiceOptions(option *IntegrationTestOption) *kernel.ServiceOptions {
	ctx := context.Background()
	//channelSecret := foundation.EnvStringOrPanic("LINE_CHANNEL_SECRET")
	//channelAccessToken := foundation.EnvStringOrPanic("LINE_ACCESS_TOKEN")
	//gcpProjectID := foundation.EnvStringOrPanic("GCP_PROJECT_ID")
	//gcpCredentials := foundation.EnvStringOrPanic("GCP_CREDENTIALS_JSON")

	gcpProjectID := "test"

	datastoreClient, err := datastore.NewLocalClient(ctx, gcpProjectID, kernel.DatastoreLocalEndpoint)
	if err != nil {
		panic(err)
	}

	pubsubClient, err := pubsub.NewLocalClient(ctx, gcpProjectID, kernel.PubsubLocalEndpoint)
	if err != nil {
		panic(err)
	}

	serviceOptions := &kernel.ServiceOptions{
		DatastoreClient: datastoreClient,
		PubsubClient:    pubsubClient,
	}

	return serviceOptions
}
