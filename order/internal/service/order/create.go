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
	if err := validatePartUUIDs(req); err != nil {
		return model.Order{}, err
	}

	partUUIDs := req.PartUUIDs()

	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return model.Order{}, fmt.Errorf(
			"получить детали из InventoryService: %w",
			err,
		)
	}

	totalPrice, err := validateParts(req, parts)
	if err != nil {
		return model.Order{}, err
	}

	now := time.Now().UTC()
	orderUUID := uuid.NewString()

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
		CreatedAt:       now,
	}

	orderItems := makeOrderItems(orderUUID, now, parts)

	err = s.txManager.Do(
		ctx,
		func(ctx context.Context) error {
			if err := s.orderRepo.Create(ctx, newOrder); err != nil {
				return fmt.Errorf(
					"сохранить заказ: %w",
					err,
				)
			}

			if err := s.orderItemRepo.Create(ctx, orderItems); err != nil {
				return fmt.Errorf(
					"сохранить детали заказа: %w",
					err,
				)
			}

			return nil
		},
	)
	if err != nil {
		return model.Order{}, err
	}

	return newOrder, nil
}

func validatePartUUIDs(req model.CreateOrderRequest) error {
	parts := []struct {
		name string
		uuid *string
	}{
		{name: "hull_uuid", uuid: &req.HullUUID},
		{name: "engine_uuid", uuid: &req.EngineUUID},
		{name: "shield_uuid", uuid: req.ShieldUUID},
		{name: "weapon_uuid", uuid: req.WeaponUUID},
	}

	for _, part := range parts {
		if part.uuid == nil {
			continue
		}

		if _, err := uuid.Parse(*part.uuid); err != nil {
			return fmt.Errorf("%s %q: %w", part.name, *part.uuid, errs.ErrInvalidUUID)
		}
	}

	return nil
}

func validateParts(req model.CreateOrderRequest, parts []model.Part) (int64, error) {
	expectedTypes, err := expectedPartTypes(req)
	if err != nil {
		return 0, err
	}

	if len(parts) != len(expectedTypes) {
		return 0, errs.ErrPartNotFound
	}

	seen := make(map[string]struct{}, len(parts))
	var totalPrice int64

	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return 0, fmt.Errorf("деталь %s: %w", part.UUID, errs.ErrOutOfStock)
		}

		expectedType, ok := expectedTypes[part.UUID]
		if !ok {
			return 0, fmt.Errorf("InventoryService вернул незапрошенную деталь %q", part.UUID)
		}

		if part.PartType != expectedType {
			return 0, fmt.Errorf(
				"деталь %q имеет тип %q, ожидается %q: %w",
				part.UUID,
				part.PartType,
				expectedType,
				errs.ErrInvalidPartType,
			)
		}

		if _, ok := seen[part.UUID]; ok {
			return 0, fmt.Errorf("InventoryService вернул дубликат детали %q", part.UUID)
		}

		seen[part.UUID] = struct{}{}
		totalPrice += part.Price
	}

	return totalPrice, nil
}

func expectedPartTypes(req model.CreateOrderRequest) (map[string]model.PartType, error) {
	expectedTypes := make(map[string]model.PartType, len(req.PartUUIDs()))

	parts := []struct {
		uuid     *string
		partType model.PartType
	}{
		{uuid: &req.HullUUID, partType: model.PartTypeHull},
		{uuid: &req.EngineUUID, partType: model.PartTypeEngine},
		{uuid: req.ShieldUUID, partType: model.PartTypeShield},
		{uuid: req.WeaponUUID, partType: model.PartTypeWeapon},
	}

	for _, part := range parts {
		if part.uuid == nil {
			continue
		}

		if previousType, ok := expectedTypes[*part.uuid]; ok && previousType != part.partType {
			return nil, fmt.Errorf(
				"деталь %q указана для типов %q и %q: %w",
				*part.uuid,
				previousType,
				part.partType,
				errs.ErrInvalidPartType,
			)
		}

		expectedTypes[*part.uuid] = part.partType
	}

	return expectedTypes, nil
}

func makeOrderItems(orderUUID string, createdAt time.Time, parts []model.Part) []model.OrderItem {
	items := make([]model.OrderItem, 0, len(parts))

	for _, part := range parts {
		items = append(items, model.OrderItem{
			UUID:      uuid.NewString(),
			OrderUUID: orderUUID,
			PartUUID:  part.UUID,
			PartType:  part.PartType,
			Price:     part.Price,
			CreatedAt: createdAt,
		})
	}

	return items
}
