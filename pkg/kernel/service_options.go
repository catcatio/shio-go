package kernel

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
)

type ServiceOptions struct {
	DatastoreClient *datastore.Client
	PubsubClient    *pubsub.Client
}

//type ChannelOptions struct {
//	ID string
//	*LineChatOptions
//	*DialogflowOptions
//	*PubsubOptions
//	*GCPOptions
//}

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
