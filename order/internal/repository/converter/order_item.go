package converter

import (
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/record"
)

func FromModelToRecordOrderItem(item model.OrderItem) record.OrderItem {
	return record.OrderItem{
		UUID:      item.UUID,
		OrderUUID: item.OrderUUID,
		PartUUID:  item.PartUUID,
		PartType:  string(item.PartType),
		Price:     item.Price,
		CreatedAt: item.CreatedAt,
	}
}

func FromRecordToModelOrderItem(item record.OrderItem) model.OrderItem {
	return model.OrderItem{
		UUID:      item.UUID,
		OrderUUID: item.OrderUUID,
		PartUUID:  item.PartUUID,
		PartType:  model.PartType(item.PartType),
		Price:     item.Price,
		CreatedAt: item.CreatedAt,
	}
}
