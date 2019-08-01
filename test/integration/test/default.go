package test

import (
	"fmt"
	"github.com/catcatio/shio-go/cmd/shio-svc/starter"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
	"os"
	"testing"
	"time"
)

type IntegrationTestOption struct {
	Name        string
	Host        string
	Port        string
	StartServer bool
}

func IntegrationTestSetup(m *testing.M, options *IntegrationTestOption, initializers ...*kernel.Initializer) {
	if testing.Short() {
		fmt.Println("skipping setup integration test")
		m.Run()
	}

	serviceOptions := makeServiceOptions(options)
	s := starter.Starter{
		Host:    options.Host,
		Port:    options.Port,
		Options: serviceOptions,
		Log:     logger.New(fmt.Sprintf("%s_test", options.Name)),
	}

	service := s.Start()

	time.Sleep(time.Second)
	c := m.Run()

	time.Sleep(time.Second)
	service.Stop()
	os.Exit(c)
}

func makeServiceOptions(option *IntegrationTestOption) *kernel.ServiceOptions {
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
