package payment

import (
	"context"
	"fmt"
	"log/slog"

	paymentErrors "github.com/Ilya96s/rocket-factory-new/payment/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/payment/internal/model"
	"github.com/google/uuid"
)

func (s *service) Pay(_ context.Context, req model.PayOrderRequest) (model.PayOrderResponse, error) {
	if _, err := uuid.Parse(req.OrderUUID); err != nil {
		return model.PayOrderResponse{}, fmt.Errorf(
			"order_uuid: %q, %w", req.OrderUUID, paymentErrors.ErrInvalidOrderUUID,
		)
	}

	switch req.PaymentMethod {
	case model.PaymentMethodCard,
		model.PaymentMethodSBP,
		model.PaymentMethodCreditCard,
		model.PaymentMethodInvestorMoney:
	default:
		return model.PayOrderResponse{}, fmt.Errorf(
			"payment_method %q: %w",
			req.PaymentMethod,
			paymentErrors.ErrInvalidPaymentMethod,
		)
	}

	slog.Info(
		"оплата прошла успешно",
		"order_uuid", req.OrderUUID,
		"payment_method", req.PaymentMethod,
	)

	return model.PayOrderResponse{
		TransactionUUID: uuid.NewString(),
	}, nil
}
