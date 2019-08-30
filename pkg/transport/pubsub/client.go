package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/catcatio/shio-go/pkg/entities/v1"
)

const (
	IncomingEventTopicName = "incoming-event-topic"
	SendMessageTopicName   = "send-message-topic"
)

type Clients struct {
	client *pubsub.Client
	topics map[string]*pubsub.Topic
}

func NewClients(client *pubsub.Client) *Clients {
	return &Clients{
		client: client,
		topics: make(map[string]*pubsub.Topic),
	}
}

func (c *Clients) topic(topicID string) *pubsub.Topic {
	if t, ok := c.topics[topicID]; ok {
		return t
	}

	t := c.client.Topic(topicID)
	c.topics[topicID] = t
	return t
}

func (c *Clients) PubsubClient() *pubsub.Client {
	return c.client
}

func (c *Clients) IncomingEventTopic() *pubsub.Topic {
	return c.topic(IncomingEventTopicName)
}

func (c *Clients) SendMessageTopic() *pubsub.Topic {
	return c.topic(SendMessageTopicName)
}

func (c *Clients) PublishSendMessageInput(ctx context.Context, input *entities.SendMessageInput) error {
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = c.SendMessageTopic().Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	return err
}

func (c *Clients) PublishIncomingEvent(ctx context.Context, input *entities.IncomingEvent) error {
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = c.IncomingEventTopic().Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	return err
}
