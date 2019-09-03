package entities

import (
	"cloud.google.com/go/datastore"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/kernel"
)

type WebhookInput struct {
	Provider  entities.ProviderType     `json:"provider"`
	ChannelID string                    `json:"channel_id"`
	Events    []*entities.IncomingEvent `json:"events"`
}

type ChannelConfig struct {
	ID                      string                     `json:"id"`
	IntentDetector          string                     `json:"intent_detector"`
	EnableLanguageDetection bool                       `json:"enable_language_detection"`
	LineChatOptions         *kernel.LineChatOptions    `json:"line_chat_options"`
	DialogflowOptions       *kernel.DialogflowOptions  `json:"dialogflow_options"`
	FulfillmentOptions      *kernel.FulfillmentOptions `json:"fulfillment_options"`
}

func (c *ChannelConfig) Load(p []datastore.Property) error {
	err := datastore.LoadStruct(c, p)
	return err
}

func (c *ChannelConfig) Save() ([]datastore.Property, error) {
	props, err := datastore.SaveStruct(c)

	for i, p := range props {
		if p.Name == "CredentialsJson" {
			props[i].NoIndex = true
		}
	}
	return props, err
}
