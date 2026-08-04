package errs

import "errors"

var (
	ErrInvalidOrderUUID     = errors.New("невалидный UUID заказа")
	ErrInvalidPaymentMethod = errors.New("невалидный способ оплаты")
)
