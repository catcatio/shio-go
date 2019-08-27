package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/catcatio/shio-go/pkg/entities/v1"
)

const (
	SendMessageTopicName = "send-message-topic"
)

type Clients struct {
	sendMessageTopic *pubsub.Topic
}

func NewClients(client *pubsub.Client) *Clients {
	return &Clients{
		sendMessageTopic: client.Topic(SendMessageTopicName),
	}
}

func (c *Clients) PublishSendMessageInput(ctx context.Context, input entities.SendMessageInput) error {
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = c.sendMessageTopic.Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	return err
}
