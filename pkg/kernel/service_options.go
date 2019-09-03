package kernel

import (
	"cloud.google.com/go/datastore"
	pubsub2 "github.com/catcatio/shio-go/pkg/transport/pubsub"
)

type ServiceOptions struct {
	DatastoreClient *datastore.Client
	PubsubClients   *pubsub2.Clients
}

type LineChatOptions struct {
	ChannelSecret      string `json:"line_channel_secret"`
	ChannelAccessToken string `json:"line_channel_access_token"`
}

type DialogflowOptions struct {
	*GCPOptions
}

type GCPOptions struct {
	ProjectID       string `json:"project_id"`
	CredentialsJson string `json:"credentials_json" datastore:",noindex"`
	Endpoint        string `json:"endpoint"`
}

type ServiceConfigBase struct {
	HttpPort string
	GrpcPort string
	Enable   bool
}

type PubsubOptions struct {
	*GCPOptions
}

type DatastoreOptions struct {
	*GCPOptions
}

type ServiceOptionsProvider interface {
	Get() (*ServiceOptions, error)
}

type FulfillmentOptions struct {
	Endpoint string `json:"fulfillment_endpoint"`
}
