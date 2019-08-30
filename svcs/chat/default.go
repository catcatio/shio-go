package chat

import (
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/transports"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/octofoxio/foundation/logger"
	"google.golang.org/grpc"
	"net/http"
)

func Initializer(options *kernel.ServiceOptions) (*grpc.Server, http.Handler) {
	chat, _ := initialUsecases(options)
	eps := endpoints.New(chat)
	httpHandler := transports.NewHttpServer(eps)

	return nil, httpHandler
}

func initialUsecases(options *kernel.ServiceOptions) (chat usecases.Chat, line usecases.Line) {
	log := logger.New("chat").WithServiceInfo("initialUsecases")
	log.Debug("enter")
	defer log.Debug("exit")

	incomingRepo := repositories.NewIncomingEventRepository(options.PubsubClients)
	channelOptions := repositories.NewChannelOptionsRepository(options.DatastoreClient)
	userProfileRepo := repositories.NewUserProfileRepository(options.DatastoreClient)
	sendMessageRepo := repositories.NewSendMessageRepository(options.PubsubClients)
	lineRepo := repositories.NewLineRepository()

	line = usecases.NewLine(userProfileRepo, lineRepo)
	chat = usecases.NewChat(incomingRepo, channelOptions, sendMessageRepo)

	return
}

func NewWebhookHandler(options *kernel.ServiceOptions) http.Handler {
	_, handler := Initializer(options)
	return handler
}

func NewPubsubHandler(options *kernel.ServiceOptions) pubsub.Handler {
	channelOptions := repositories.NewChannelOptionsRepository(options.DatastoreClient)
	sendMessageRepo := repositories.NewSendMessageRepository(options.PubsubClients)
	sendMessageUsecase := usecases.NewSendMessageUsecase(channelOptions, sendMessageRepo)
	intentRepo := repositories.NewIntentRepository(channelOptions)
	intentUsecase := usecases.NewIntentUsecase(channelOptions, intentRepo)

	handlers := make(endpoints.PubsubMessageHandlers)
	handlers[pubsub.SendMessageTopicName] = endpoints.NewSendMessagePubsubEndpoint(sendMessageUsecase)
	handlers[pubsub.IncomingEventTopicName] = endpoints.NewIncomingEventPubsubEndpoint(intentUsecase)

	return transports.NewPubsubServer(handlers)
}
