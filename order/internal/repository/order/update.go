package order

import (
	"context"
	"fmt"

	errs "github.com/Ilya96s/rocket-factory-new/order/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	const query = `
	UPDATE orders
	SET status = $1, transaction_uuid = $2, payment_method = $3, updated_at = NOW()
	WHERE uuid = $4`

	cmdTag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		order.Status,
		order.TransactionUUID,
		order.PaymentMethod,
		order.UUID,
	)
	if err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}
