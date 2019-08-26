package repositories

import (
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineRepository interface {
	GetUserProfile(ctx context.Context, channelSecret, channelAccessToken, userID string) (*entities.UserProfile, error)
}

type lineRepository struct {
}

func NewLineRepository() LineRepository {
	return &lineRepository{}
}

func (l *lineRepository) GetUserProfile(ctx context.Context, channelSecret, channelAccessToken, userID string) (*entities.UserProfile, error) {
	// TODO: use caching
	if lineClient, err := linebot.New(channelSecret, channelAccessToken); err != nil {
		return nil, err
	} else {
		lineProfile, err := lineClient.GetProfile(userID).WithContext(ctx).Do()
		profile := new(entities.UserProfile)
		if err == nil {
			profile.ID = lineProfile.UserID
			profile.DisplayName = lineProfile.DisplayName
			profile.PictureUrl = lineProfile.PictureURL
		} else {
			profile.ID = userID
		}

		return profile, err
	}
}
