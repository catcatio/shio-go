package main

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/internal/prime"
	"github.com/catcatio/shio-go/nub/datastore"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/octofoxio/foundation"
	"strings"
)

func main() {
	ctx := context.Background()
	endpoint := "localhost:5545"

	channelSecret := foundation.EnvStringOrPanic("LINE_CHANNEL_SECRET")
	channelAccessToken := foundation.EnvStringOrPanic("LINE_ACCESS_TOKEN")
	gcpProjectID := foundation.EnvStringOrPanic("GCP_PROJECT_ID")
	gcpCredentials := foundation.EnvStringOrPanic("GCP_CREDENTIALS_JSON")
	c := strings.Trim(gcpCredentials, "\"")
	c = strings.Replace(c, `'`, `"`, -1)
	//datastoreClient, err := datastore.NewClientWithCredentials(ctx, gcpProjectID, c)
	datastoreClient, err := datastore.NewLocalClient(ctx, gcpProjectID, endpoint)
	chConfig := &entities.ChannelConfig{
		ID: "too-bright",
		LineChatOptions: &kernel.LineChatOptions{
			ChannelSecret:      channelSecret,
			ChannelAccessToken: channelAccessToken,
		},
		DialogflowOptions: &kernel.DialogflowOptions{
			GCPOptions: &kernel.GCPOptions{
				ProjectID:       gcpProjectID,
				CredentialsJson: gcpCredentials,
			},
		},
	}

	fmt.Println(err)
	key, err := prime.ChannelConfig(datastoreClient, chConfig)

	if err != nil {
		panic(err)
	}

	cfg := new(entities.ChannelConfig)

	err = datastoreClient.Get(ctx, key, cfg)

	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

}
