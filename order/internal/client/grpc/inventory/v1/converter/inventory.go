package converter

import (
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
)

func PartsToModel(parts []*inventoryv1.Part) []model.Part {
	outParts := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		outParts = append(outParts, partToModel(part))
	}

	return outParts
}

func partToModel(part *inventoryv1.Part) model.Part {
	return model.Part{
		UUID:          part.GetUuid(),
		Name:          part.GetName(),
		PartType:      partTypeToModel(part.GetPartType()),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
	}
}

func partTypeToModel(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	case inventoryv1.PartType_PART_TYPE_UNSPECIFIED:
		return model.PartTypeUnspecified
	default:
		return model.PartTypeUnspecified
	}
}
