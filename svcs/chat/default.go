package chat

import (
	"fmt"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/catcatio/shio-go/svcs/chat/transports"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"github.com/gorilla/mux"
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

func writeResponse(w http.ResponseWriter, status int, msg interface{}) {
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "%s: %v", http.StatusText(status), msg)
}

func New(options *kernel.ServiceOptions) http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(fmt.Sprintf("/chat/{%s}/{%s}", "provider", "channleid"), func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		writeResponse(w, 200, params)
	}).
		Methods("POST")
	muxRouter.HandleFunc("/_ping", func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, 200, "pong")
	}).
		Methods("GET")
	muxRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, 200, "yae yae")
	}).
		Methods("GET")

	return muxRouter
}
