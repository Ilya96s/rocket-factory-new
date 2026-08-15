package part

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/converter"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/record"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	const query = `SELECT
    uuid,
    name,
    description,
    part_type,
    price,
    stock_quantity,
    created_at
	FROM parts
	WHERE uuid = $1`

	var part record.Part

	err := r.pool.QueryRow(ctx, query, uuid).Scan(
		&part.UUID,
		&part.Name,
		&part.Description,
		&part.PartType,
		&part.Price,
		&part.StockQuantity,
		&part.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Part{}, errs.ErrPartNotFound
		}

		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return converter.FromRecordToModel(part), nil
}
