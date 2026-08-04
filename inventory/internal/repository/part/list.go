package part

import (
	"context"

	errs "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/converter"
)

func (r *repository) List(_ context.Context, filter model.PartFilter) ([]model.Part, error) {
	var parts []model.Part
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Если передан список partUUIDs
	if len(filter.UUIDs) > 0 {
		for _, uuid := range filter.UUIDs {
			part, ok := r.parts[uuid]
			if !ok {
				return []model.Part{}, errs.ErrPartNotFound
			}
			parts = append(parts, converter.FromRecordToModel(part))
		}
	}

	// Если тип не указан (PART_TYPE_UNSPECIFIED)
	if filter.PartType == model.PartTypeUnspecified {
		for _, part := range r.parts {
			parts = append(parts, converter.FromRecordToModel(part))
		}
		return parts, nil
	}

	// Фильтрация по типу
	for _, part := range r.parts {
		if part.PartType == string(filter.PartType) {
			parts = append(parts, converter.FromRecordToModel(part))
		}
	}

	return parts, nil
}
