package order

import (
	"context"
	"fmt"
	"time"

	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	if _, err := uuid.Parse(req.HullUUID); err != nil {
		return model.Order{}, fmt.Errorf("hull_uuid %q: %w", req.HullUUID, errs.ErrInvalidUUID)
	}
	if _, err := uuid.Parse(req.EngineUUID); err != nil {
		return model.Order{}, fmt.Errorf("engine_uuid %q: %w", req.EngineUUID, errs.ErrInvalidUUID)
	}
	if req.ShieldUUID != nil {
		if _, err := uuid.Parse(*req.ShieldUUID); err != nil {
			return model.Order{}, fmt.Errorf("shield_uuid %q: %w", req.ShieldUUID, errs.ErrInvalidUUID)
		}
	}
	if req.WeaponUUID != nil {
		if _, err := uuid.Parse(*req.WeaponUUID); err != nil {
			return model.Order{}, fmt.Errorf("weapon_uuid %q: %w", req.WeaponUUID, errs.ErrInvalidUUID)
		}
	}

	partUUIDs := req.PartUUIDs()

	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить деталь из InventoryService: %w", err)
	}

	if len(parts) != len(partUUIDs) {
		return model.Order{}, errs.ErrPartNotFound
	}

	var totalPrice int64
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.UUID, errs.ErrOutOfStock)
		}
		totalPrice += part.Price
	}

	orderUUID := uuid.New().String()

	newOrder := model.Order{
		UUID:            orderUUID,
		HullUUID:        req.HullUUID,
		EngineUUID:      req.EngineUUID,
		ShieldUUID:      req.ShieldUUID,
		WeaponUUID:      req.WeaponUUID,
		TotalPrice:      totalPrice,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          model.OrderStatusPendingPayment,
		CreatedAt:       time.Now().UTC(),
	}
	if err = s.orderRepo.Create(ctx, newOrder); err != nil {
		return model.Order{}, fmt.Errorf("сохранить заказ: %w", err)
	}

	return newOrder, nil
}
