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

	chat := usecases.NewIncomingEventUsecase(channelConfig, pubsubRepo)
	eps := endpoints.New(chat)
	httpHandler := transports.NewHttpServer(eps)
	return httpHandler
}

func NewPubsubHandler(options *kernel.ServiceOptions) pubsub.Handler {
	channelConfig := repositories.NewChannelConfigRepository(options.DatastoreClient)
	pubsubRepo := repositories.NewPubsubChannelRepository(options.PubsubClients)
	outgoingEventUsecase := usecases.NewOutgoingEventUsecase(channelConfig, pubsubRepo)
	intentRepo := repositories.NewIntentRepository(channelConfig)
	intentUsecase := usecases.NewIntentUsecase(channelConfig, intentRepo, pubsubRepo)
	fulfillmentRepo := repositories.NewFulfillmentRepository(channelConfig)
	fulfillmentUsecase := usecases.NewFulfillmentUsecase(fulfillmentRepo)

	handlers := make(endpoints.PubsubMessageHandlers)
	handlers[pubsub.OutgoingEventTopicName] = endpoints.NewOutgoingEventPubsubEndpoint(outgoingEventUsecase)
	handlers[pubsub.IncomingEventTopicName] = endpoints.NewIncomingEventPubsubEndpoint(intentUsecase)
	handlers[pubsub.FulfillmentTopicName] = endpoints.NewFulfillmentPubsubEndpoint(fulfillmentUsecase)

	return transports.NewPubsubServer(handlers)
}
