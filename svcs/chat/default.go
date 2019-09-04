package chat

import (
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/transports"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"net/http"
)

func NewWebhookHandler(options *kernel.ServiceOptions) http.Handler {
	pubsubRepo := repositories.NewPubsubChannelRepository(options.PubsubClients)
	channelConfig := repositories.NewChannelConfigRepository(options.DatastoreClient)

	chat := usecases.NewIncomingEventUsecase(pubsubRepo)
	eps := endpoints.New(chat, channelConfig, pubsubRepo)
	httpHandler := transports.NewHttpServer(eps)
	return httpHandler
}

func NewPubsubHandler(options *kernel.ServiceOptions) pubsub.Handler {
	channelConfig := repositories.NewChannelConfigRepository(options.DatastoreClient)
	pubsubRepo := repositories.NewPubsubChannelRepository(options.PubsubClients)
	outgoingRepo := repositories.NewOutgoingRepo()
	outgoingEventUsecase := usecases.NewOutgoingEventUsecase(outgoingRepo)
	intentRepo := repositories.NewIntentRepository()
	intentUsecase := usecases.NewIntentUsecase(intentRepo, pubsubRepo)
	fulfillmentRepo := repositories.NewFulfillmentRepository()
	fulfillmentUsecase := usecases.NewFulfillmentUsecase(fulfillmentRepo)

	handlers := make(endpoints.PubsubMessageHandlers)
	handlers[pubsub.OutgoingEventTopicName] = endpoints.NewOutgoingEventPubsubEndpoint(outgoingEventUsecase, channelConfig)
	handlers[pubsub.IncomingEventTopicName] = endpoints.NewIncomingEventPubsubEndpoint(intentUsecase, channelConfig)
	handlers[pubsub.FulfillmentTopicName] = endpoints.NewFulfillmentPubsubEndpoint(fulfillmentUsecase, channelConfig)

	return transports.NewPubsubServer(handlers)
}
