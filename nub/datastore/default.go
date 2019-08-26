package datastore

import (
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func NewLocalClient(ctx context.Context, projectID string, apiEndpoint string) (*datastore.Client, error) {
	return datastore.NewClient(ctx, projectID,
		option.WithEndpoint(apiEndpoint),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
	)
}

func NewClientWithCredentials(ctx context.Context, projectID string, credentialsJson string) (*datastore.Client, error) {
	return datastore.NewClient(ctx, projectID,
		option.WithCredentialsJSON([]byte(credentialsJson)),
	)
}

func Newlient(ctx context.Context, projectID string) (*datastore.Client, error) {
	return datastore.NewClient(ctx, projectID)
}

func NewDefaultClient(ctx context.Context, projectID string) (*datastore.Client, error) {
	return datastore.NewClient(ctx, projectID)
}
