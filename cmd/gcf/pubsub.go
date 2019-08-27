package gcf

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat"
	"github.com/octofoxio/foundation/logger"
	"os"
)

var pubsubHandler pubsub.Handler

func init() {
	projectID := os.Getenv("GCLOUD_PROJECT")
	ctx := context.Background()
	fmt.Printf("initial function pubsub in %s\n", projectID)
	options := newServiceOptions(ctx, projectID)
	pubsubHandler = chat.NewPubsubHandler(options)
}

func handlePubsub(topic string, ctx context.Context, m pubsub.RawPubsubMessage) (err error) {
	if pubsubHandler == nil {
		err := fmt.Errorf("handler is nil")
		logger.Default.WithServiceInfo("HandlePubsub").WithError(err).Errorf("failed to handler message")
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if e := r.(error); e != nil {
				err = e
			}
			logger.Default.WithServiceInfo("HandlePubsub").Errorf("something wrong: %v", r)
		}
	}()

	logger.Default.Printf("pubsub message: %v", m.Data)
	logger.Default.Printf("pubsub message: (string) %v", string(m.Data))

	return pubsubHandler.Serve(topic, ctx, m)
}

func HandleSendMessagePubsub(ctx context.Context, m pubsub.RawPubsubMessage) error {
	return handlePubsub(pubsub.SendMessageTopicName, ctx, m)
}
