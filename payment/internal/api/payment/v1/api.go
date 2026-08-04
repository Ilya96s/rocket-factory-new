package v1

import paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"

type handler struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

func New(paymentService PaymentService) *handler {
	return &handler{
		paymentService: paymentService,
	}
}
