package endpoints

import (
	"context"
	entities2 "github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"net/http"
)

type EndpointFunc func(w http.ResponseWriter, r *http.Request)
type ProviderEndpointHandlers map[string]ProviderEndpointFunc
type ProviderEndpointFunc func(ctx context.Context, options *entities.ChannelConfig, w http.ResponseWriter, r *http.Request)

type Endpoints struct {
	Webhook  EndpointFunc
	Outgoing EndpointFunc
	Ping     EndpointFunc
}

var (
	ParamChannelID = "channelid"
	ParamProvider  = "provider"
	ParamAction    = "action"
)

func New(incomingEvent usecases.IncomingEventUsecase, channelConfigRepo repositories.ChannelConfigRepository, pubsub repositories.PubsubChannelRepository) *Endpoints {
	handlers := make(ProviderEndpointHandlers)
	handlers[entities2.ProviderTypeLine.String()] = newLineEndpointFunc(incomingEvent)

	return &Endpoints{
		Webhook:  newWebhookEndpoint(channelConfigRepo, handlers),
		Outgoing: newOutgoingEndpoint(channelConfigRepo, pubsub),
		Ping:     newPingEndpoint(),
	}
}
