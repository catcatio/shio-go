package entities

import (
	"cloud.google.com/go/datastore"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/pkg/kernel"
)

type WebhookInput struct {
	Provider  entities.ProviderType
	ChannelID string
	Events    []*entities.IncomingEvent `json:"events"`
}

type ChannelConfig struct {
	ID             string
	IntentDetector string
	DetectLanguage bool
	*kernel.LineChatOptions
	*kernel.DialogflowOptions
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
