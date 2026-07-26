package order

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Order, error) {
	return s.orderRepo.Get(ctx, uuid)
}
