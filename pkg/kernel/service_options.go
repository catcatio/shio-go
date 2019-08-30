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
	ChannelSecret      string
	ChannelAccessToken string
}

type DialogflowOptions struct {
	*GCPOptions
}

type GCPOptions struct {
	ProjectID       string
	CredentialsJson string
	Endpoint        string
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
