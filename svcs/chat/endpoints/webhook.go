package endpoints

import (
	"github.com/catcatio/shio-go/nub"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/gorilla/mux"
	"github.com/octofoxio/foundation/logger"
	"net/http"
)

func newWebhookEndpoint(channelConfigRepo repositories.ChannelConfigRepository, handlers ProviderEndpointHandlers) EndpointFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := nub.NewID()
		log := logger.New("WebhookEndpoint").WithServiceInfo("handler").WithRequestID(requestID)
		log.Debug("enter")
		defer log.Debug("exit")

		params := mux.Vars(r)
		channelID := params[ParamChannelID]
		provider := params[ParamProvider]

		log = log.WithField("channelID", channelID).WithField("provider", provider)

		if channelID == "" || provider == "" {
			log.Error("parameter missing")
			writeNotFoundResponse(w)
			return
		}

		handler := handlers[provider]
		if handler == nil {
			log.Error("handler not found")
			log.Error("parameter missing")
			writeNotFoundResponse(w)
			return
		}

		ctx := shio.NewContextWithRequestID(requestID)
		log.Info("getting channel config")
		channelConfig, err := channelConfigRepo.Get(ctx, channelID)
		if err != nil {
			log.WithError(err).Error("get channel config failed")
			writeNotFoundResponse(w)
			return
		}

		ctx = shio.PutContextValue(ctx, "channel_config", *channelConfig)

		log.Info("handling request config")
		handler(ctx, channelConfig, w, r)
	}
}
