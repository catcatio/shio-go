package shio

import (
	"context"
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/logger"
	"strconv"
)

func EnvIntOrPanic(key string) (port int) {
	portStr := foundation.EnvStringOrPanic(key)

	if value, err := strconv.ParseInt(portStr, 10, 32); err == nil {
		port = int(value)
	} else {
		logger.New("envIntOrPanic").Panic(err)
	}

	return
}

func EnvStringOrPanic(key string) string {
	return foundation.EnvStringOrPanic(key)
}

func EnvInt(key string, defaultValue int) (port int) {
	portStr := foundation.EnvString(key, fmt.Sprintf("%d", defaultValue))

	port = defaultValue
	if value, err := strconv.ParseInt(portStr, 10, 32); err == nil {
		port = int(value)
	}

	return
}

func EnvInt64(key string, defaultValue int64) (port int64) {
	portStr := foundation.EnvString(key, fmt.Sprintf("%d", defaultValue))

	if value, err := strconv.ParseInt(portStr, 10, 32); err == nil {
		port = value
	}

	return
}

func EnvFloat64(key string, defaultValue float64) (port float64) {
	portStr := foundation.EnvString(key, fmt.Sprintf("%.7f", defaultValue))

	if value, err := strconv.ParseFloat(portStr, 64); err == nil {
		port = value
	}

	return
}

func EnvString(key string, defaultValue string) string {
	return foundation.EnvString(key, defaultValue)
}

func NewContext(ctx context.Context) context.Context {
	requestID := foundation.GetRequestIDFromContext(ctx)
	return NewContextWithRequestID(requestID)
}

func NewContextWithRequestID(requestID string) context.Context {
	newCtx := context.Background()
	if requestID != "" {
		return AppendRequestIDToContext(newCtx, requestID)
	}

	return newCtx
}

func AppendRequestIDToContext(ctx context.Context, userID string) context.Context {
	ctx = context.WithValue(ctx, foundation.FoundationRequestIDContextKey, userID)
	ctx = context.WithValue(ctx, foundation.GRPC_METADATA_REQUEST_ID_KEY, userID)
	ctx = foundation.AppendRequestIDToContext(ctx, userID)
	return ctx
}
