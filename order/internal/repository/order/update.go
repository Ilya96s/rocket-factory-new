package order

import (
	"context"

	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/converter"
)

func (r *repository) Update(_ context.Context, order model.Order) error {
	orderRecord := converter.FromModelToRecord(order)
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.UUID]; !ok {
		return errs.ErrOrderNotFound
	}
	r.orders[order.UUID] = orderRecord
	return nil
}
