package order

import (
	"context"
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return model.Order{}, fmt.Errorf(
			"получить заказ %q: %w",
			uuid,
			err,
		)
	}

	items, err := s.orderItemRepo.GetByOrderUUID(ctx, uuid)
	if err != nil {
		return model.Order{}, fmt.Errorf(
			"получить позиции заказа %q: %w",
			uuid,
			err,
		)
	}

	for _, item := range items {
		switch item.PartType {
		case model.PartTypeHull:
			order.HullUUID = item.PartUUID

		case model.PartTypeEngine:
			order.EngineUUID = item.PartUUID

		case model.PartTypeShield:
			partUUID := item.PartUUID
			order.ShieldUUID = &partUUID

		case model.PartTypeWeapon:
			partUUID := item.PartUUID
			order.WeaponUUID = &partUUID
		}
	}

	if order.HullUUID == "" {
		return model.Order{}, fmt.Errorf(
			"заказ %q не содержит корпус",
			uuid,
		)
	}

	if order.EngineUUID == "" {
		return model.Order{}, fmt.Errorf(
			"заказ %q не содержит двигатель",
			uuid,
		)
	}

	return order, nil
}
