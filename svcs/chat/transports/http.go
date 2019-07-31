package transports

import (
	"github.com/catcatio/shio-go/svcs/chat/endpoints"
	"github.com/gorilla/mux"
	"net/http"
)

func NewHttpServer(eps *endpoints.Endpoints) http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/chat/line", eps.LineWebhook).
		Methods("POST")
	m.HandleFunc("/_ping", eps.Ping).
		Methods("GET")

	return m
}
