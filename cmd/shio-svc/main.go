package main

import (
	"github.com/catcatio/shio-go/cmd/shio-svc/starter"
	"github.com/catcatio/shio-go/pkg/kernel"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
	"os"
	"os/signal"
	"syscall"
)

func loadOptions(platform string) *kernel.ServiceOptions {
	switch platform {
	default:
		return starter.LoadLocalOptions()
	}
}

type stopper interface {
	Stop()
}

func main() {
	log := logger.New("shio-svc").WithServiceInfo("main")
	log.Println("hurray !!!")

	platform := foundation.EnvString("PLATFORM", "local")
	httpHost := foundation.EnvString("HOST", "localhost")
	httpPort := foundation.EnvString("PORT", "30001")

	log.Printf("using platform: %s", platform)
	serviceOptions := loadOptions(platform)

	s := &starter.Starter{
		Host:    httpHost,
		Port:    httpPort,
		Options: serviceOptions,
		Log:     log,
	}

	service := s.Start()
	// gracefully stop service
	gracefullyStop(log, service)
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
