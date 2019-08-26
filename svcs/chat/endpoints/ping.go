package endpoints

import (
	"fmt"
	"net/http"
)

func newPingEndpoint() EndpointFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "_pong")
	}
}
