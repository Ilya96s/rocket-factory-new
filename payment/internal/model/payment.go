package model

type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
)

type PayOrderRequest struct {
	OrderUUID     string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID string
}
