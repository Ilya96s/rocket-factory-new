package converter

import (
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/record"
)

func FromRecordToModel(rec record.PartRecord) model.Part {
	return model.Part{
		UUID:          rec.UUID,
		Name:          rec.Name,
		Description:   rec.Description,
		Price:         rec.Price,
		PartType:      toPartType(rec.PartType),
		StockQuantity: rec.StockQuantity,
		CreatedAt:     rec.CreatedAt,
	}
}

func toPartType(partType string) model.PartType {
	switch partType {
	case "HULL":
		return model.PartTypeHull
	case "ENGINE":
		return model.PartTypeEngine
	case "WEAPON":
		return model.PartTypeWeapon
	case "SHIELD":
		return model.PartTypeShield
	default:
		return model.PartTypeUnspecified
	}
}
