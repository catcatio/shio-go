package starter

import (
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/octofoxio/foundation"
)

func LoadLocalOptions() *kernel.ServiceOptions {
	channelSecret := foundation.EnvStringOrPanic("LINE_CHANNEL_SECRET")
	channelAccessToken := foundation.EnvStringOrPanic("LINE_ACCESS_TOKEN")
	gcpProjectName := foundation.EnvStringOrPanic("GCP_PROJECT_NAME")
	gcpCredential := foundation.EnvStringOrPanic("GCP_CREDENTIALS")

	serviceOptions := &kernel.ServiceOptions{
		LineChatOptions: &kernel.LineChatOptions{
			ChannelSecret:      channelSecret,
			ChannelAccessToken: channelAccessToken,
		},
		DialogflowOptions: &kernel.DialogflowOptions{
			GCPOptions: &kernel.GCPOptions{
				ProjectName:    gcpProjectName,
				CredentialJson: gcpCredential,
			},
		},
	}

	return serviceOptions
}
