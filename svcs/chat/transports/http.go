package transports

import (
	"fmt"
	"github.com/catcatio/shio-go/pkg/transport/middleware"
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/gorilla/mux"
	"net/http"
)

func NewHttpServer(eps *endpoints.Endpoints) http.Handler {
	m := mux.NewRouter()
	m.Use(middleware.DoneWriterMiddleware)
	m.HandleFunc(fmt.Sprintf("/chat/{%s}/{%s}", endpoints.ParamProvider, endpoints.ParamChannelID), eps.Webhook).
		Methods("POST")
	m.HandleFunc("/_ping", eps.Ping).
		Methods("GET")

	return m
}
