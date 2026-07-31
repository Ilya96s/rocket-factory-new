package converter

import (
	"fmt"

	errs "github.com/Ilya96s/rocket-factory-new/payment/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/payment/internal/model"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToModel(req *paymentv1.PayOrderRequest) (model.PayOrderRequest, error) {
	if req == nil {
		return model.PayOrderRequest{}, fmt.Errorf("пустой запрос: %w", errs.ErrInvalidOrderUUID)
	}
	payMethod := paymentMethodToModel(req.GetPaymentMethod())

	return model.PayOrderRequest{
		OrderUUID:     req.GetOrderUuid(),
		PaymentMethod: payMethod,
	}, nil
}

func paymentMethodToModel(method paymentv1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func PayOrderResultToProto(result model.PayOrderResponse) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: result.TransactionUUID,
	}
}
