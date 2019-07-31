package endpoints

import (
	"fmt"
	"github.com/catcatio/shio-go/svcs/chat/usecases"
	"net/http"
)

type EndpointFunc func(w http.ResponseWriter, r *http.Request)

type Endpoints struct {
	LineWebhook EndpointFunc
	Ping        EndpointFunc
}

func NewPingEndpoint() EndpointFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "_pong")
	}
}

func New(line usecases.Line) *Endpoints {
	return &Endpoints{
		LineWebhook: NewLineWebhookEndpointFunc(line),
		Ping:        NewPingEndpoint(),
	}
}
