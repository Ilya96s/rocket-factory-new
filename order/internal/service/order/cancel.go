package order

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	foundOrder, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return err // TODO какая ошибка
	}

	if foundOrder.Status == model.OrderStatusPaid || foundOrder.Status == model.OrderStatusCanceled {
		return nil // TODO какая ошибка
	}

	foundOrder.Status = model.OrderStatusPendingPayment
	if err = s.orderRepo.Update(ctx, foundOrder); err != nil {
		return nil // TODO какая ошибка ?
	}

	return nil
}
