package gcf

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/svcs/chat"
	"github.com/octofoxio/foundation/logger"
	"net/http"
	"os"
)

var webHookHandler http.Handler

func init() {
	projectID := os.Getenv("GCLOUD_PROJECT")
	ctx := context.Background()
	fmt.Printf("initial function webhook handler in %s\n", projectID)
	options := newServiceOptions(ctx, projectID)
	webHookHandler = chat.NewWebhookHandler(options)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if webHookHandler == nil {
		err := fmt.Errorf("handler is nil")
		logger.Default.WithServiceInfo("HandleRequest").WithError(err).Errorf("failed to handler request")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Default.WithServiceInfo("HandleRequest").Errorf("something wrong: %v", r)
			_, _ = fmt.Fprintf(w, "%s: %v", http.StatusText(http.StatusInternalServerError), r)
		}
	}()

	webHookHandler.ServeHTTP(w, r)
}
