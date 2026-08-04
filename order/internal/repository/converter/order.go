package converter // TODO какое название пакета правильнее использовать converter или order ?
import (
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/record"
)

// FromModelToRecord - конвертирует модель сервисного слоя в модель хранения данных.
func FromModelToRecord(order model.Order) record.OrderRecord {
	orderRecord := record.OrderRecord{
		UUID:       order.UUID,
		HulUUID:    order.HullUUID,
		EngineUUID: order.EngineUUID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
	}
	if order.ShieldUUID != nil {
		orderRecord.ShieldUUID = *order.ShieldUUID
	}
	if order.WeaponUUID != nil {
		orderRecord.WeaponUUID = *order.WeaponUUID
	}
	if order.PaymentMethod != nil {
		orderRecord.PaymentMethod = string(*order.PaymentMethod)
	}
	orderRecord.Status = string(order.Status)

	return orderRecord
}

// FromRecordToModel - конвертирует модель хранения данных в модель сервисного слоя.
func FromRecordToModel(orderRecord record.OrderRecord) (model.Order, error) {
	order := model.Order{
		UUID:       orderRecord.UUID,
		HullUUID:   orderRecord.HulUUID,
		EngineUUID: orderRecord.EngineUUID,
		TotalPrice: orderRecord.TotalPrice,
		CreatedAt:  orderRecord.CreatedAt,
	}
	if orderRecord.ShieldUUID != "" {
		order.ShieldUUID = &orderRecord.ShieldUUID
	}
	if orderRecord.WeaponUUID != "" {
		order.WeaponUUID = &orderRecord.WeaponUUID
	}
	if orderRecord.PaymentMethod != "" {
		method, err := toPaymentMethod(orderRecord.PaymentMethod)
		if err != nil {
			return model.Order{}, err // TODO какая ошибка ?
		}
		order.PaymentMethod = &method
	}
	if orderRecord.Status != "" {
		status, err := toStatus(orderRecord.Status)
		if err != nil {
			return model.Order{}, err // TODO какая ошибка ?
		}
		order.Status = status
	}
	return order, nil
}

// toPaymentMethod - конвертирует строковый тип оплаты в тип PaymentMethod
func toPaymentMethod(method string) (model.PaymentMethod, error) {
	switch method {
	case "CARD":
		return model.PaymentMethodCard, nil
	case "SBP":
		return model.PaymentMethodSBP, nil
	case "CREDIT_CARD":
		return model.PaymentMethodCreditCard, nil
	case "INVESTOR_MONEY":
		return model.PaymentMethodInvestorMoney, nil
	default:
		return "", fmt.Errorf("неизвестный способ оплаты: %s", method)
	}
}

// toStatus - конвертирует строковый тип статуса в тип OrderStatus
func toStatus(status string) (model.OrderStatus, error) {
	switch status {
	case "PENDING_PAYMENT":
		return model.OrderStatusPendingPayment, nil
	case "PAID":
		return model.OrderStatusPaid, nil
	case "CANCELLED":
		return model.OrderStatusCanceled, nil
	default:
		return "", fmt.Errorf("неизвестный статус: %s", status)
	}
}
