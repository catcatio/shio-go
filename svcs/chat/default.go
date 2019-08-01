package chat

import (
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/transports"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	lineUtil "github.com/catcatio/shio-go/utils/line"
	"google.golang.org/grpc"
	"net/http"
)

func Initializer(serviceOptions *kernel.ServiceOptions) (*grpc.Server, http.Handler) {
	lineOptions := serviceOptions.LineChatOptions
	dialogflowOptions := serviceOptions.DialogflowOptions

	parser := lineUtil.NewParser(lineOptions.ChannelSecret)
	intentRepo := repositories.NewDialogflowIntentRepository(dialogflowOptions.ProjectName, dialogflowOptions.CredentialJson)

	line := usecases.New(lineOptions, parser, intentRepo)
	eps := endpoints.New(line)
	httpHandler := transports.NewHttpServer(eps)

	return nil, httpHandler
}
