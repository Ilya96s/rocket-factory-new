package v1

import (
	"context"
	"net/http"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	orderv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func (h *handler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	serviceRequest := model.CreateOrderRequest{
		HullUUID:   req.GetHullUUID().String(),
		EngineUUID: req.GetEngineUUID().String(),
	}

	if shieldUUID, ok := req.ShieldUUID.Get(); ok {
		value := shieldUUID.String()
		serviceRequest.ShieldUUID = &value
	}
	if weaponUUID, ok := req.WeaponUUID.Get(); ok {
		value := weaponUUID.String()
		serviceRequest.WeaponUUID = &value
	}

	order, err := h.orderService.Create(ctx, model.CreateOrderRequest{})
	if err != nil {
		return mapCreateError(err), nil
	}

	orderUUID, err := uuid.Parse(order.UUID)
	if err != nil {
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "order service вернул невалидный order_uuid",
		}, nil
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
