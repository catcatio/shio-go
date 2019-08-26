package bot_webhook_fn

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/nub/datastore"
	"github.com/catcatio/shio-go/nub/pubsub"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/svcs/chat"
	"github.com/octofoxio/foundation/logger"
	"net/http"
	"os"
)

var handler http.Handler

func init() {
	projectID := os.Getenv("GCLOUD_PROJECT")
	ctx := context.Background()
	fmt.Printf("initial router: %s\n", projectID)

	pubsubClient, err := pubsub.NewDefaultClient(ctx, projectID)
	if err != nil {
		logger.Default.Panicf("create pubsub client failed %s", err.Error())
	}

	datastoreClient, err := datastore.NewDefaultClient(ctx, projectID)
	if err != nil {
		logger.Default.Panicf("create datastore client failed %s", err.Error())
	}

	options := &kernel.ServiceOptions{
		PubsubClient:    pubsubClient,
		DatastoreClient: datastoreClient,
	}

	handler = chat.New(options)
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		logger.Default.WithServiceInfo("Webhook").Error("handler is nil")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Default.WithServiceInfo("Webhook").Errorf("panic in webhook: %v", r)
			_, _ = fmt.Fprintf(w, "%s: %v", http.StatusText(http.StatusInternalServerError), r)
		}
	}()

	handler.ServeHTTP(w, r)
}
