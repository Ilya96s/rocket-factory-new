package part

import (
	"context"

	errs "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return converter.FromRecordToModel(part), nil
}
