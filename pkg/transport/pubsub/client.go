package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/catcatio/shio-go/pkg/entities/v1"
)

const (
	IncomingEventTopicName = "incoming-event-topic"
	OutgoingEventTopicName = "outgoing-event-topic"
	FulfillmentTopicName   = "fulfillment-topic"
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

func (c *Clients) OutgoingEventTopic() *pubsub.Topic {
	return c.topic(OutgoingEventTopicName)
}

func (c *Clients) FulfillmentTopic() *pubsub.Topic {
	return c.topic(FulfillmentTopicName)
}

func (c *Clients) publish(ctx context.Context, topic *pubsub.Topic, message interface{}) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = topic.Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	return err
}

func (c *Clients) PublishOutgoingEventInput(ctx context.Context, input *entities.OutgoingEvent) error {
	return c.publish(ctx, c.OutgoingEventTopic(), input)
}

func (c *Clients) PublishIncomingEvent(ctx context.Context, input *entities.IncomingEvent) error {
	return c.publish(ctx, c.IncomingEventTopic(), input)
}

func (c *Clients) PublishFulfillmentInput(ctx context.Context, input *entities.IncomingEvent) error {
	return c.publish(ctx, c.FulfillmentTopic(), input)
}
