package order

import (
	"context"

	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return converter.FromRecordToModel(order)
}
