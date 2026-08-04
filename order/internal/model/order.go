package model

import "time"

// Order - внутренняя модель заказа для работы на уровне сервисного слоя.
type Order struct {
	UUID            string
	HullUUID        string
	EngineUUID      string
	ShieldUUID      *string
	WeaponUUID      *string
	TotalPrice      int64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// PaymentMethod - внутренний способ оплаты.
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// OrderStatus - внутренний статус заказа.
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCanceled       OrderStatus = "CANCELLED"
)

// CreateOrderRequest - внутренняя модель запроса.
type CreateOrderRequest struct {
	HullUUID   string
	EngineUUID string
	ShieldUUID *string
	WeaponUUID *string
}

func (r *CreateOrderRequest) PartUUIDs() []string {
	uuids := []string{r.HullUUID, r.EngineUUID}
	if r.ShieldUUID != nil {
		uuids = append(uuids, *r.ShieldUUID)
	}
	if r.WeaponUUID != nil {
		uuids = append(uuids, *r.WeaponUUID)
	}
	return uuids
}
