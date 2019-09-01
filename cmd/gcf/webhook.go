package gcf

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/transport/middleware"
	"github.com/catcatio/shio-go/svcs/chat"
	"github.com/octofoxio/foundation/logger"
	"net/http"
	"os"
)

var webHookHandler http.Handler

func init() {
	projectID := os.Getenv("GCLOUD_PROJECT")
	log := logger.New("handleWebhook").WithServiceInfo("init")
	log.Printf("initial function webhook handler in %s\n", projectID)

	ctx := context.Background()
	options := newServiceOptions(ctx, projectID)
	webHookHandler = chat.NewWebhookHandler(options)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	log := logger.New("handleWebhook").WithServiceInfo("handleRequest")
	if webHookHandler == nil {
		err := fmt.Errorf("handler is nil")
		log.WithError(err).Errorf("failed to handler request")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("something wrong: %v", r)
			_, _ = fmt.Fprintf(w, "%s: %v", http.StatusText(http.StatusInternalServerError), r)
		}
	}()

	middleware.DoneWriterMiddleware(webHookHandler).ServeHTTP(w, r)
}
