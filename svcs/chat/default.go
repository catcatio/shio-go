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
	incomingRepo := repositories.NewIncomingEventRepository(options.PubsubClients)
	channelConfig := repositories.NewChannelConfigRepository(options.DatastoreClient)

	chat := usecases.NewIncomingEventUsecase(channelConfig, incomingRepo)
	eps := endpoints.New(chat)
	httpHandler := transports.NewHttpServer(eps)
	return httpHandler
}

func NewPubsubHandler(options *kernel.ServiceOptions) pubsub.Handler {
	channelConfig := repositories.NewChannelConfigRepository(options.DatastoreClient)
	sendMessageRepo := repositories.NewSendMessageRepository(options.PubsubClients)
	sendMessageUsecase := usecases.NewSendMessageUsecase(channelConfig, sendMessageRepo)
	intentRepo := repositories.NewIntentRepository(channelConfig)
	intentUsecase := usecases.NewIntentUsecase(channelConfig, intentRepo)

	handlers := make(endpoints.PubsubMessageHandlers)
	handlers[pubsub.SendMessageTopicName] = endpoints.NewSendMessagePubsubEndpoint(sendMessageUsecase)
	handlers[pubsub.IncomingEventTopicName] = endpoints.NewIncomingEventPubsubEndpoint(intentUsecase)

	return transports.NewPubsubServer(handlers)
}
