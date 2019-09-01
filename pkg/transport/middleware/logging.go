package middleware

import (
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

func LoggingHandler(logger *logger.Logger, h http.Handler) http.Handler {
	return handlers.LoggingHandler(&localLogger{log: logger}, h)
}
