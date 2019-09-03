package main

import (
	datastore2 "cloud.google.com/go/datastore"
	pubsub3 "cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/catcatio/shio-go/nub/datastore"
	"github.com/catcatio/shio-go/nub/pubsub"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/pkg/transport/middleware"
	pubsub2 "github.com/catcatio/shio-go/pkg/transport/pubsub"
	"github.com/catcatio/shio-go/svcs/chat"
	"github.com/catcatio/shio-go/svcs/chat/entities"
	"github.com/gorilla/mux"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func newServiceOptions(ctx context.Context, projectID string) *kernel.ServiceOptions {
	log := logger.New("local").WithServiceInfo("newServiceOptions")
	datastoreEndpoint := shio.EnvString("GCP_DATASTORE_ENDPOINT", "localhost:5545")
	pubsubEndpoint := shio.EnvString("GCP_PUBSUB_ENDPOINT", "localhost:8085")

	datastoreClient, err := datastore.NewLocalClient(ctx, projectID, datastoreEndpoint)
	if err != nil {
		log.Panic(err)
	}

	pubsubClient, err := pubsub.NewLocalClient(ctx, projectID, pubsubEndpoint)

	if err != nil {
		log.Panic(err)
	}

	return &kernel.ServiceOptions{
		DatastoreClient: datastoreClient,
		PubsubClients:   pubsub2.NewClients(pubsubClient),
	}
}

func newPubsubPushEndpoint(handler pubsub2.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		topic := params["topic"]
		msg := &struct {
			Message struct {
				Attributes map[string]string
				Data       []byte
				ID         string `json:"message_id"`
			}
			Subscription string
		}{}

		if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
			http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
			return
		}

		rawMessage := pubsub2.RawPubsubMessage{
			Data: msg.Message.Data,
		}

		if e := handler.Serve(topic, r.Context(), rawMessage); e != nil {
			fmt.Println(e)
			// http.Error(w, fmt.Sprintf("failed to handle messge: %v", e), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func createPubsubChannel(opts *kernel.ServiceOptions, topicName string) {
	ctx := context.Background()
	topic := opts.PubsubClients.PubsubClient().Topic(topicName)
	if exist, _ := topic.Exists(ctx); !exist {
		_, err := opts.PubsubClients.PubsubClient().CreateTopic(ctx, topicName)
		if err != nil {
			panic(err)
		}
	}

	subscriptionName := topicName + "-subscription"
	subscriptions := opts.PubsubClients.PubsubClient().Subscription(subscriptionName)

	if exist, err := subscriptions.Exists(ctx); exist && err == nil {
		_ = subscriptions.Delete(ctx)
	}

	_, err := opts.PubsubClients.PubsubClient().CreateSubscription(ctx, subscriptionName, pubsub3.SubscriptionConfig{
		Topic: topic,
		PushConfig: pubsub3.PushConfig{
			Endpoint: "http://host.docker.internal:30001/pubsub/" + topicName,
		},
		AckDeadline:         0,
		RetainAckedMessages: false,
		RetentionDuration:   0,
		ExpirationPolicy:    nil,
		Labels:              nil,
	})

	if err != nil {
		panic(err)
	}
}

func setChannelConfig(opts *kernel.ServiceOptions, config *entities.ChannelConfig) {
	ctx := context.Background()
	key := datastore2.NameKey("ChannelConfig", "no-boundary", nil)
	_, err := opts.DatastoreClient.Put(ctx, key, config)
	if err != nil {
		panic(err)
	}
}

func setupEnvironment(opts *kernel.ServiceOptions) {
	createPubsubChannel(opts, pubsub2.IncomingEventTopicName)
	createPubsubChannel(opts, pubsub2.OutgoingEventTopicName)
	createPubsubChannel(opts, pubsub2.FulfillmentTopicName)

	channelSecret := foundation.EnvStringOrPanic("LINE_CHANNEL_SECRET")
	channelAccessToken := foundation.EnvStringOrPanic("LINE_ACCESS_TOKEN")
	gcpProjectID := foundation.EnvStringOrPanic("GCP_PROJECT_ID")
	gcpCredentials := foundation.EnvStringOrPanic("GCP_CREDENTIALS_JSON")

	setChannelConfig(opts, &entities.ChannelConfig{
		ID:             "no-boundary",
		IntentDetector: "dialogflow",
		LineChatOptions: &kernel.LineChatOptions{
			ChannelSecret:      channelSecret,
			ChannelAccessToken: channelAccessToken,
		},
		DialogflowOptions: &kernel.DialogflowOptions{
			GCPOptions: &kernel.GCPOptions{
				ProjectID:       gcpProjectID,
				CredentialsJson: gcpCredentials,
				Endpoint:        "",
			},
		},
		FulfillmentOptions: &kernel.FulfillmentOptions{
			Endpoint: foundation.EnvStringOrPanic("FULFILLMENT_ENDPOINT"),
		},
	})
}

func makeEndpoints(httpHandler http.Handler, pubsubHttpHandler func(http.ResponseWriter, *http.Request)) http.Handler {
	var m *mux.Router

	if mm, ok := httpHandler.(*mux.Router); ok {
		m = mm
	} else {
		panic("not *mux.Router")
	}

	m.HandleFunc("/pubsub/{topic}", pubsubHttpHandler)

	return m
}

func main() {
	log := logger.New("shio-local").WithServiceInfo("main")
	log.Println("hurray !!!")

	httpHost := foundation.EnvString("HOST", "localhost")
	httpPort := foundation.EnvString("PORT", "30001")
	projectID := "shio-local"

	serviceOptions := newServiceOptions(context.Background(), projectID)
	setupEnvironment(serviceOptions)

	s := &Starter{
		Host:    httpHost,
		Port:    httpPort,
		Options: serviceOptions,
		Log:     log,
	}

	service := s.Start()
	// gracefully stop service
	gracefullyStop(log, service)
}

type Service struct {
	httpServer *http.Server
	log        *logger.Logger
}

func (s *Service) Stop() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		s.log.WithError(err).Error("http server shut down failed")
	}
}

type Starter struct {
	Host    string
	Port    string
	Options *kernel.ServiceOptions
	Log     *logger.Logger
}

type localLogger struct {
	log *logger.Logger
}

func (l *localLogger) Write(p []byte) (n int, err error) {
	l.log.Info(strings.TrimSpace(string(p)))
	return len(p), nil
}

func (s *Starter) Start() *Service {
	if s.Log == nil {
		s.Log = logger.New("shio-svc")
	}
	log := s.Log.WithServiceInfo("starter")
	log.Println("starting...")

	webHookHandler := chat.NewWebhookHandler(s.Options)
	chat.NewPubsubHandler(s.Options)
	pubsubHandler := chat.NewPubsubHandler(s.Options)
	pubsubHttpHandler := newPubsubPushEndpoint(pubsubHandler)

	handler := makeEndpoints(webHookHandler, pubsubHttpHandler)
	finalHandler := middleware.LoggingHandler(log.WithServiceInfo("serve"), handler)
	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: finalHandler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.Log.WithError(err).Fatalf("http server error")
		}
	}()

	log.Printf("started on %s", address)
	return &Service{
		httpServer: srv,
		log:        s.Log,
	}
}

type stopper interface {
	Stop()
}

func gracefullyStop(log *logger.Logger, stopper stopper) {
	gracefullyClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGINT)
		sig := <-sigint

		log.Printf("🤷 🤷 ‍system signal received: %s", sig.String())
		// stop service
		log.Info("stopping services")
		stopper.Stop()
		close(gracefullyClosed)
	}()
	<-gracefullyClosed
	log.Info("🙋 🙋 bye")
}
