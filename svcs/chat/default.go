package chat

import (
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/transports"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/octofoxio/foundation/logger"
	"google.golang.org/grpc"
	"net/http"
)

func Initializer(options *kernel.ServiceOptions) (*grpc.Server, http.Handler) {
	log := logger.New("chatsvc").WithServiceInfo("Initializer")
	log.Debug("enter")
	defer log.Debug("exit")

	incomingRepo := repositories.NewIncomingEventRepository(options.PubsubClient)
	channelOptions := repositories.NewChannelOptionsRepository(options.DatastoreClient)
	userProfileRepo := repositories.NewUserProfileRepository(options.DatastoreClient)
	lineRepo := repositories.NewLineRepository()

	line := usecases.NewLine(userProfileRepo, lineRepo)
	chat := usecases.NewChat(incomingRepo, channelOptions)
	eps := endpoints.New(chat, line)
	httpHandler := transports.NewHttpServer(eps)

	return nil, httpHandler
}
