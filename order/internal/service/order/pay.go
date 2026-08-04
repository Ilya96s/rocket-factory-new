package order

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

func (s *service) Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error) {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return "", err // TODO какая ошибка ?
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", nil // TODO какая ошибка ?
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, uuid, method)
	if err != nil {
		return "", err // TODO какая ошибка ?
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method

	if err = s.orderRepo.Update(ctx, order); err != nil {
		return "", err // TODO какая ошибка ?
	}

	return transactionUUID, nil
}
