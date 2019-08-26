package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func NewLocalClient(ctx context.Context, projectID string, apiEndpoint string) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectID,
		option.WithEndpoint(apiEndpoint),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
	)
}

func NewClientWithCredentials(ctx context.Context, projectID string, credentialsJson string) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectID,
		option.WithCredentialsJSON([]byte(credentialsJson)),
	)
}

func NewClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectID)
}
