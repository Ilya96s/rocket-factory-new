package converter

import (
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
)

func PartsToModel(parts []*inventoryv1.Part) []model.Part {
	outParts := make([]model.Part, len(parts))

	for _, part := range parts {
		outParts = append(outParts, partToModel(part))
	}

	return outParts
}

func partToModel(part *inventoryv1.Part) model.Part {
	return model.Part{
		UUID:          part.GetUuid(),
		Name:          part.GetName(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
	}
}
