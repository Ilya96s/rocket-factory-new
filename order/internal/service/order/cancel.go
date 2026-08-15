package order

import (
	"context"
	"fmt"

	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	foundOrder, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return fmt.Errorf(
			"получить заказ %q: %w",
			uuid,
			err,
		)
	}

	switch foundOrder.Status {
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCanceled:
		return errs.ErrOrderCanceled
	case model.OrderStatusPendingPayment:
		// Заказ можно отменить.
	default:
		return fmt.Errorf("заказ %q имеет неизвестный статус %q", uuid, foundOrder.Status)
	}

	foundOrder.Status = model.OrderStatusCanceled
	if err = s.orderRepo.Update(ctx, foundOrder); err != nil {
		return fmt.Errorf(
			"отменить заказ %q: %w",
			uuid,
			err,
		)
	}

	return nil
}
