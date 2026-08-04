package v1

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, req model.PayOrderRequest) (model.PayOrderResponse, error)
}
