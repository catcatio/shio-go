package usecases

import (
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/octofoxio/foundation/logger"
)

type Line interface {
	GetUserProfile(ctx context.Context, channelSecret, channelAccessToken, userID string) (*entities.UserProfile, error)
}

type line struct {
	provider        entities.ProviderType
	userProfileRepo repositories.UserProfileRepository
	lineRepo        repositories.LineRepository
	log             *logger.Logger
}

func NewLine(userProfileRepo repositories.UserProfileRepository, lineRepo repositories.LineRepository) Line {
	return &line{
		provider:        entities.ProviderTypeLine,
		userProfileRepo: userProfileRepo,
		lineRepo:        lineRepo,
		log:             logger.New("line"),
	}
}

func (l *line) GetUserProfile(ctx context.Context, channelSecret, channelAccessToken, userID string) (profile *entities.UserProfile, err error) {
	// TODO: use caching
	profile, err = l.lineRepo.GetUserProfile(ctx, channelSecret, channelAccessToken, userID)

	if err != nil {
		l.log.WithError(err).Println("get user profile failed")
		return
	}

	err = l.userProfileRepo.Put(ctx, "line", profile)

	if err != nil {
		l.log.WithError(err).Println("save user profile failed")
		return
	}

	return
}
