package order

import (
	"context"
	"fmt"

	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

func (s *service) Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error) {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return "", fmt.Errorf("получить заказ %q: %w", uuid, err)
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return "", errs.ErrOrderAlreadyPaid
	case model.OrderStatusCanceled:
		return "", errs.ErrOrderCanceled
	case model.OrderStatusPendingPayment:
		// Заказ можно оплатить.
	default:
		return "", fmt.Errorf("заказ %q имеет неизвестный статус %q", uuid, order.Status)
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, uuid, method)
	if err != nil {
		return "", fmt.Errorf("оплатить заказ %q: %w", uuid, err)
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method

	if err = s.orderRepo.Update(ctx, order); err != nil {
		return "", fmt.Errorf("сохранить оплату заказа %q: %w", uuid, err)
	}

	return transactionUUID, nil
}
