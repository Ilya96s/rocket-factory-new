package order

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order model.Order) error {
	orderRecord := converter.FromModelToRecord(order)

	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.UUID] = orderRecord
	return nil
}
