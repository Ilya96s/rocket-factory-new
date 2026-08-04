package part

import (
	"context"

	errs "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, partUUID string) (model.Part, error) {
	UUID, err := uuid.Parse(partUUID)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}
	return s.partRepository.Get(ctx, UUID.String())
}
