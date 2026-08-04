package v1

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/payment/internal/converter"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	modelRequest, err := converter.PayOrderRequestToModel(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	payOrderResponse, err := h.paymentService.Pay(ctx, modelRequest)
	if err != nil {
		return nil, mapError(err)
	}

	return converter.PayOrderResultToProto(payOrderResponse), nil
}
