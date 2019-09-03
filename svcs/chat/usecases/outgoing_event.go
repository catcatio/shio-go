package usecases

import (
	"context"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/octofoxio/foundation/logger"
)

type OutgoingEventUsecase interface {
	Handle(ctx context.Context, input *entities.OutgoingEvent) error
}

type outgoingEventUsecase struct {
	outgoingRepo repositories.OutgoingRepository
	log          *logger.Logger
}

func (s *outgoingEventUsecase) Handle(ctx context.Context, input *entities.OutgoingEvent) error {
	log := s.log.WithServiceInfo("Handle")

	err := s.outgoingRepo.Handler(ctx, input)
	if err != nil {
		log.WithError(err).WithField("channelID", input.ChannelID).Error("handle event failed")
		return err
	}

	return nil

}

func NewOutgoingEventUsecase(outgoingRepo repositories.OutgoingRepository) OutgoingEventUsecase {
	return &outgoingEventUsecase{
		outgoingRepo: outgoingRepo,
		log:          logger.New("OutgoingEventUsecase"),
	}
}
