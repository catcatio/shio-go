package endpoints

import (
	"context"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

func newWebHookEndpoint(chat usecases.Chat, handlers ProviderEndpointHandlers) EndpointFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		channelID := params[ParamChannelID]
		provider := params[ParamProvider]
		handler := handlers[provider]

		if channelID == "" || provider == "" || handler == nil {
			writeNotFoundResponse(w)
			return
		}

		ctx := context.Background()
		channelOptions, err := chat.GetChannelConfig(ctx, channelID)

		if err != nil {
			writeNotFoundResponse(w)
			return
		}

		handler(ctx, channelOptions, w, r)
	}
}
