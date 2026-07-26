package v1

import (
	"context"
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/service/order"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ order.InventoryClient = (*client)(nil)

// client - обертка над proto client
type client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

func New(c inventoryv1.InventoryServiceClient) *client {
	return &client{
		inventoryClient: c,
	}
}

func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.inventoryClient.ListParts(
		ctx,
		&inventoryv1.ListPartsRequest{
			Uuids: uuids,
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, fmt.Errorf(
				"получить список деталей: %w",
				err,
			)
		}

		switch st.Code() {
		case codes.InvalidArgument:
			return nil, fmt.Errorf(
				"inventory service: %w",
				errs.ErrInvalidUUID,
			)

		case codes.NotFound:
			return nil, fmt.Errorf(
				"inventory service: %w",
				errs.ErrPartNotFound,
			)

		default:
			return nil, fmt.Errorf(
				"inventory service, код %s: %w",
				st.Code(),
				err,
			)
		}
	}

	return converter.PartsToModel(resp.GetParts()), nil
}
