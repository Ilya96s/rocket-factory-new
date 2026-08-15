package v1

import (
	"context"
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/service/order"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ order.PaymentClient = (*client)(nil)

// client - обертка над proto client
type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

func New(p paymentv1.PaymentServiceClient) *client {
	return &client{
		paymentClient: p,
	}
}

func (c *client) PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	resp, err := c.paymentClient.PayOrder(
		ctx,
		&paymentv1.PayOrderRequest{
			OrderUuid:     orderUUID,
			PaymentMethod: converter.PaymentMethodToProto(method),
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", fmt.Errorf(
				"оплатить заказ: %w",
				err,
			)
		}

		switch st.Code() {
		case codes.InvalidArgument:
			return "", fmt.Errorf(
				"payment service: %w",
				errs.ErrInvalidUUID,
			)

		default:
			return "", fmt.Errorf(
				"payment service, код %s: %w",
				st.Code(),
				err,
			)
		}
	}

	return resp.GetTransactionUuid(), nil
}
