package endpoints

import (
	"context"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"net/http"
)

type EndpointFunc func(w http.ResponseWriter, r *http.Request)
type ProviderEndpointHandlers map[string]ProviderEndpointFunc
type ProviderEndpointFunc func(ctx context.Context, options *entities.ChannelConfig, w http.ResponseWriter, r *http.Request)

type Endpoints struct {
	Webhook EndpointFunc
	Ping    EndpointFunc
}

var (
	ParamChannelID = "channelid"
	ParamProvider  = "provider"
)

func New(chat usecases.Chat, line usecases.Line) *Endpoints {
	handlers := ProviderEndpointHandlers{
		"line": newLineEndpointFunc(chat, line),
	}

	return &Endpoints{
		Webhook: newWebHookEndpoint(chat, handlers),
		Ping:    newPingEndpoint(),
	}
}
