package usecases

import (
	"context"
	"fmt"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
)

type IntentUsecase interface {
	HandleEvents(ctx context.Context, in *entities.IncomingEvent) (err error)
}

type intentUsecase struct {
	channelOptionsRepo repositories.ChannelOptionsRepository
	intentRepo         repositories.IntentRepository
}

func NewIntentUsecase(channelOptionsRepo repositories.ChannelOptionsRepository, intentRepo repositories.IntentRepository) IntentUsecase {
	return &intentUsecase{channelOptionsRepo: channelOptionsRepo, intentRepo: intentRepo}
}

func (i *intentUsecase) HandleEvents(ctx context.Context, in *entities.IncomingEvent) (err error) {
	fmt.Println(in)
	// TODO: implement me
	return nil
}
