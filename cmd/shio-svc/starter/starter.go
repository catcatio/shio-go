package starter

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/catcatio/shio-go/svcs/chat"
	"github.com/gorilla/handlers"
	"github.com/octofoxio/foundation/logger"
	"net/http"
	"strings"
)

type localLogger struct {
	log *logger.Logger
}

func (l *localLogger) Write(p []byte) (n int, err error) {
	l.log.Info(strings.TrimSpace(string(p)))
	return len(p), nil
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

func (s *Starter) Start() *Service {
	if s.Log == nil {
		s.Log = logger.New("shio-svc")
	}
	log := s.Log.WithServiceInfo("starter")
	log.Println("starting...")

	_, httpHandler := chat.Initializer(s.Options)
	loggedRouter := handlers.LoggingHandler(&localLogger{log: s.Log.WithServiceInfo("http-server")}, httpHandler)

	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: loggedRouter,
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
