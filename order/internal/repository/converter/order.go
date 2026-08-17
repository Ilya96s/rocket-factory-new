package converter

import (
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/record"
)

func FromModelToRecordOrder(modelOrder model.Order) record.Order {
	var paymentMethod *string
	if modelOrder.PaymentMethod != nil {
		pm := string(*modelOrder.PaymentMethod)
		paymentMethod = &pm
	}

	return record.Order{
		UUID:            modelOrder.UUID,
		TotalPrice:      modelOrder.TotalPrice,
		Status:          string(modelOrder.Status),
		TransactionUUID: modelOrder.TransactionUUID,
		PaymentMethod:   paymentMethod,
		CreatedAt:       modelOrder.CreatedAt,
	}
}

func FromRecordToModelOrder(recordOrder record.Order) model.Order {
	order := model.Order{
		UUID:            recordOrder.UUID,
		TotalPrice:      recordOrder.TotalPrice,
		TransactionUUID: recordOrder.TransactionUUID,
		CreatedAt:       recordOrder.CreatedAt,
	}

	if recordOrder.PaymentMethod != nil {
		method := model.PaymentMethod(*recordOrder.PaymentMethod)
		order.PaymentMethod = &method
	}

	order.Status = model.OrderStatus(recordOrder.Status)

	return order
}
