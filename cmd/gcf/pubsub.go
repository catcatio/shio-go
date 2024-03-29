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
	log := logger.New("handlePubsub").WithServiceInfo("init")
	log.Printf("initial function pubsub in %s\n", projectID)

	ctx := context.Background()
	options := newServiceOptions(ctx, projectID)
	pubsubHandler = chat.NewPubsubHandler(options)
}

func handlePubsub(topic string, ctx context.Context, m pubsub.RawPubsubMessage) (err error) {
	log := logger.New("handlePubsub").WithServiceInfo("handlePubsub").WithField("topic", topic)
	if m.Data == nil {
		err := fmt.Errorf("data is nil")
		log.WithError(err).Error("handle pubsub failed")
		return err
	}

	if pubsubHandler == nil {
		err := fmt.Errorf("handler is nil")
		log.WithServiceInfo("HandlePubsub").WithError(err).Errorf("failed to handler message")
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if e := r.(error); e != nil {
				err = e
			}
			log.WithServiceInfo("HandlePubsub").Errorf("something wrong: %v", r)
		}
	}()

	log.Printf("pubsub message: %v", string(m.Data))

	return pubsubHandler.Serve(topic, ctx, m)
}

func HandleOutgoingEventPubsub(ctx context.Context, m pubsub.RawPubsubMessage) error {
	return handlePubsub(pubsub.OutgoingEventTopicName, ctx, m)
}

func HandleIncomingEventPubsub(ctx context.Context, m pubsub.RawPubsubMessage) error {
	return handlePubsub(pubsub.IncomingEventTopicName, ctx, m)
}

func HandleFulfillmentPubsub(ctx context.Context, m pubsub.RawPubsubMessage) error {
	return handlePubsub(pubsub.FulfillmentTopicName, ctx, m)
}
