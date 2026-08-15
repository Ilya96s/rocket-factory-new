package order

import (
	"context"
	"errors"
	"fmt"

	orderErrs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/converter"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/record"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Order, error) {
	const query = `
	SELECT uuid, total_price, status, transaction_uuid, payment_method, created_at, updated_at
	FROM orders where uuid = $1`

	var order record.Order
	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, uuid).Scan(
		&order.UUID,
		&order.TotalPrice,
		&order.Status,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, orderErrs.ErrOrderNotFound
		}

		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	return converter.FromRecordToModelOrder(order), nil
}
