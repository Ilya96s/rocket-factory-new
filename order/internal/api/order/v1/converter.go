package v1

import (
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	orderv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func orderToDTO(order model.Order) (*orderv1.OrderDto, error) {
	orderUUID, err := uuid.Parse(order.UUID)
	if err != nil {
		return nil, fmt.Errorf("невалидный UUID заказа: %w", err)
	}
	hullUUID, err := uuid.Parse(order.HullUUID)
	if err != nil {
		return nil, fmt.Errorf("невалидный UUID корпуса: %w", err)
	}
	engineUUID, err := uuid.Parse(order.EngineUUID)
	if err != nil {
		return nil, fmt.Errorf("невалидный UUID двигателя: %w", err)
	}
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		value, err := uuid.Parse(*order.ShieldUUID)
		if err != nil {
			return nil, fmt.Errorf("невалидный UUID щита: %w", err)
		}
		shieldUUID = orderv1.NewOptNilUUID(value)
	}
	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		value, err := uuid.Parse(*order.WeaponUUID)
		if err != nil {
			return nil, fmt.Errorf("невалидный UUID оружия: %w", err)
		}

		weaponUUID = orderv1.NewOptNilUUID(value)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		value, err := uuid.Parse(*order.TransactionUUID)
		if err != nil {
			return nil, fmt.Errorf(
				"невалидный UUID транзакции: %w",
				err,
			)
		}

		transactionUUID = orderv1.NewOptNilUUID(value)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		method, err := paymentMethodToDTO(*order.PaymentMethod)
		if err != nil {
			return nil, err
		}

		paymentMethod = orderv1.NewOptNilPaymentMethod(method)
	}

	status, err := orderStatusToDTO(order.Status)
	if err != nil {
		return nil, err
	}

	return &orderv1.OrderDto{
		OrderUUID:       orderUUID,
		HullUUID:        hullUUID,
		EngineUUID:      engineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
		CreatedAt:       order.CreatedAt,
	}, nil
}

func paymentMethodToModel(method orderv1.PaymentMethod) (model.PaymentMethod, error) {
	switch method {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard, nil

	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP, nil

	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard, nil

	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, nil

	default:
		return "", fmt.Errorf("неподдерживаемый способ оплаты: %s", method)
	}
}

func paymentMethodToDTO(method model.PaymentMethod) (orderv1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCard:
		return orderv1.PaymentMethodCARD, nil

	case model.PaymentMethodSBP:
		return orderv1.PaymentMethodSBP, nil

	case model.PaymentMethodCreditCard:
		return orderv1.PaymentMethodCREDITCARD, nil

	case model.PaymentMethodInvestorMoney:
		return orderv1.PaymentMethodINVESTORMONEY, nil

	default:
		return "", fmt.Errorf(
			"неизвестный способ оплаты: %s",
			method,
		)
	}
}

func orderStatusToDTO(status model.OrderStatus) (orderv1.OrderStatus, error) {
	switch status {
	case model.OrderStatusPendingPayment:
		return orderv1.OrderStatusPENDINGPAYMENT, nil

	case model.OrderStatusPaid:
		return orderv1.OrderStatusPAID, nil

	case model.OrderStatusCanceled:
		return orderv1.OrderStatusCANCELLED, nil

	default:
		return "", fmt.Errorf(
			"неизвестный статус заказа: %s",
			status,
		)
	}
}
