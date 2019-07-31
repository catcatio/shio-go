package main

import (
	"fmt"
	"github.com/catcatio/shio-go/svcs/chat"
	"github.com/gorilla/handlers"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
	"net/http"
)

func main() {
	log := logger.New("shio-svc")
	log.Println("hurray !!!")
	httpHost := foundation.EnvString("HOST", "localhost")
	httpPort := foundation.EnvString("PORT", "3001")
	channelSecret := foundation.EnvStringOrPanic("LINE_CHANNEL_SECRET")
	channelAccessToken := foundation.EnvStringOrPanic("LINE_ACCESS_TOKEN")

	_, httpHandler := chat.Initializer(channelSecret, channelAccessToken)

	loggedRouter := handlers.LoggingHandler(&localLogger{log: log}, httpHandler)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", httpHost, httpPort), loggedRouter)
	if err != nil {
		panic(err)
	}
}

type localLogger struct {
	log *logger.Logger
}

func (l *localLogger) Write(p []byte) (n int, err error) {
	l.log.Info(string(p))
	return len(p), nil
}
