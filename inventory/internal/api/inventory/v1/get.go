package v1

import (
	"context"
	"errors"

	"github.com/Ilya96s/rocket-factory-new/inventory/internal/converter"
	partErrors "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	part, err := h.partService.Get(ctx, req.GetUuid())
	if err != nil {
		switch {
		case errors.Is(err, partErrors.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, partErrors.ErrPartNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "внутренняя ошибка сервиса")
		}
	}

	return converter.GetPartResultToProto(part), nil
}
