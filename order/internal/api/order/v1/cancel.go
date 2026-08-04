package v1

import (
	"context"

	orderv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/openapi/order/v1"
)

func (h *handler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := h.orderService.Cancel(ctx, params.OrderUUID.String())
	if err != nil {
		return mapCancelError(err), nil
	}

	return &orderv1.CancelOrderResponse{}, nil
}
