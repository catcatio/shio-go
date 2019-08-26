package repositories

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
)

type UserProfileRepository interface {
	Get(ctx context.Context, provider, userID string) (*entities.UserProfile, error)
	Put(ctx context.Context, provider string, profile *entities.UserProfile) error
}

type userProfileRepository struct {
	dataStoreClient *datastore.Client
}

func NewUserProfileRepository(dataStoreClient *datastore.Client) UserProfileRepository {
	return &userProfileRepository{
		dataStoreClient: dataStoreClient,
	}
}

func (u *userProfileRepository) makeKey(provider, userID string) *datastore.Key {
	return datastore.NameKey("UserProfile", fmt.Sprintf("%s-%s", provider, userID), nil)
}
func (u *userProfileRepository) Get(ctx context.Context, provider, userID string) (*entities.UserProfile, error) {
	profile := new(entities.UserProfile)
	err := u.dataStoreClient.Get(ctx, u.makeKey(provider, userID), profile)

	if err == nil {
		return nil, err
	}

	return profile, nil
}

func (u *userProfileRepository) Put(ctx context.Context, provider string, profile *entities.UserProfile) error {
	if profile == nil {
		return fmt.Errorf("profile is nil")
	}

	_, err := u.dataStoreClient.Put(ctx, u.makeKey(provider, profile.ID), profile)

	if err == nil {
		return err
	}

	return nil
}
