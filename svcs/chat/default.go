package chat

import (
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/catcatio/shio-go/svcs/chat/transports"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	lineUtil "github.com/catcatio/shio-go/utils/line"
	"google.golang.org/grpc"
	"net/http"
)

func Initializer(channelSecret, channelAccessToken string) (*grpc.Server, http.Handler) {
	parser := lineUtil.NewParser(channelSecret)
	//intentRepo := repositories.NewDialogflowIntentRepository("", "{}")
	line := usecases.New(channelSecret, channelAccessToken, parser, nil)
	eps := endpoints.New(line)
	httpHandler := transports.NewHttpServer(eps)

	return nil, httpHandler
}
