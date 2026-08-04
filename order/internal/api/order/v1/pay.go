package v1

import (
	"context"
	"net/http"

	orderv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func (h *handler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	paymentMethod, err := paymentMethodToModel(req.PaymentMethod)
	if err != nil {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	transactionUUIDValue, err := h.orderService.Pay(ctx, params.OrderUUID.String(), paymentMethod)
	if err != nil {
		return mapPayError(err), nil
	}

	transactionUUID, err := uuid.Parse(transactionUUIDValue)
	if err != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "payment service вернул невалидный UUID транзакции",
		}, nil
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
