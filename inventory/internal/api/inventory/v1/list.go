package v1

import (
	"context"
	"errors"

	"github.com/Ilya96s/rocket-factory-new/inventory/internal/converter"
	inventoryErrors "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	partFilter := converter.ListPartsRequestToModel(req)
	parts, err := h.partService.List(ctx, partFilter)
	if err != nil {
		switch {
		case errors.Is(err, inventoryErrors.ErrInvalidUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, inventoryErrors.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "внутренняя ошибка сервиса")
		}
	}
	return converter.InventoryResultToProto(parts), nil
}
