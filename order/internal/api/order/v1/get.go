package v1

import (
	"context"
	"net/http"

	orderv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/openapi/order/v1"
)

func (h *handler) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := h.orderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		return mapGetError(err), nil
	}

	dto, err := orderToDTO(order)
	if err != nil {
		return &orderv1.GetOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "не удалось сформировать ответ",
		}, nil
	}

	return dto, nil
}
