package order

import (
	"context"
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	orderRecord := converter.FromModelToRecordOrder(order)

	const query = `
	INSERT INTO orders (uuid, total_price, status, created_at)
	VALUES ($1, $2, $3, $4)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		orderRecord.UUID,
		orderRecord.TotalPrice,
		orderRecord.Status,
		orderRecord.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}
