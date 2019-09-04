package endpoints

import (
	"encoding/json"
	"fmt"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/gorilla/mux"
	"github.com/octofoxio/foundation/logger"
	"net/http"
)

func newOutgoingEndpoint(channelConfigRepo repositories.ChannelConfigRepository, pubsub repositories.PubsubChannelRepository) EndpointFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.New("OutgoingEndpoint").WithServiceInfo("handler")
		log.Debug("enter")
		defer func() { log.Debug("exit") }()

		params := mux.Vars(r)
		channelID := params[ParamChannelID]

		input := new(entities.OutgoingEvent)
		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
			return
		}
		log = log.WithRequestID(input.RequestID)

		if channelID != input.ChannelID {
			log.Errorf("parameter mismatch channelID, expecting: %s, got: %s", channelID, input.ChannelID)
			writeBadRequest(w)
			return
		}

		ctx := shio.NewContextWithRequestID(input.RequestID)
		log.Infof("getting channel config %s", input.ChannelID)
		channelConfig, err := channelConfigRepo.Get(ctx, input.ChannelID)
		if err != nil || channelConfig == nil {
			log.WithError(err).Error("get channel config failed")
			writeNotFoundResponse(w)
			return
		}

		err = pubsub.PublishOutgoingEventInput(ctx, input)
		if err != nil {
			log.WithError(err).Println("publish outgoing event failed")
		}

		log.Println("message forwarded")
	}
}
