package v1

import (
	"errors"

	paymentErrors "github.com/Ilya96s/rocket-factory-new/payment/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, paymentErrors.ErrInvalidOrderUUID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, paymentErrors.ErrInvalidPaymentMethod):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "внутренняя ошибка сервиса")
	}
}
