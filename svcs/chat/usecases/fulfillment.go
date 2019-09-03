package usecases

import (
	"context"
	shio "github.com/catcatio/shio-go/pkg"
	"github.com/catcatio/shio-go/pkg/entities/v1"
	"github.com/catcatio/shio-go/svcs/chat/repositories"
	"github.com/octofoxio/foundation/logger"
)

type FulfillmentUsecase interface {
	Dispatch(ctx context.Context, input *entities.IncomingEvent) error
}

type fulfillmentUsecase struct {
	fulfillmentRepo repositories.FulfillmentRepository
	log             *logger.Logger
}

func NewFulfillmentUsecase(fulfillmentRepo repositories.FulfillmentRepository) FulfillmentUsecase {
	return &fulfillmentUsecase{
		fulfillmentRepo: fulfillmentRepo,
		log:             logger.New("FulfillmentUsecase"),
	}
}

func (f *fulfillmentUsecase) Dispatch(ctx context.Context, input *entities.IncomingEvent) error {
	log := f.log.WithServiceInfo("Dispatch").WithRequestID(shio.ReqIDFromContext(ctx))
	result, e := f.fulfillmentRepo.ResolveFulfillment(ctx, input)
	log.Println(result)

	return e
}
