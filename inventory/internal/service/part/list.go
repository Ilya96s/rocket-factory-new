package part

import (
	"context"
	"fmt"

	errs "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	"github.com/google/uuid"
)

func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) > 0 {
		for _, partUUID := range filter.UUIDs {
			if _, err := uuid.Parse(partUUID); err != nil {
				return []model.Part{}, fmt.Errorf("получить список деталей: %q, %w", partUUID, errs.ErrInvalidUUID)
			}
		}
	}

	return s.partRepository.List(ctx, filter)
}
