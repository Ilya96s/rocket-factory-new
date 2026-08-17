package order_item

import (
	"context"
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/converter"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/record"
)

func (r *repository) GetByOrderUUID(ctx context.Context, orderUUID string) ([]model.OrderItem, error) {
	const query = `
		select
			uuid,
			order_uuid,
			part_uuid,
			part_type,
			price,
			created_at
		from order_items
		where order_uuid = $1
	`

	executor := r.getter.DefaultTrOrDB(ctx, r.pool)

	rows, err := executor.Query(
		ctx,
		query,
		orderUUID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"получить позиции заказа %q: %w",
			orderUUID,
			err,
		)
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0)

	for rows.Next() {
		var itemRecord record.OrderItem

		if err := rows.Scan(
			&itemRecord.UUID,
			&itemRecord.OrderUUID,
			&itemRecord.PartUUID,
			&itemRecord.PartType,
			&itemRecord.Price,
			&itemRecord.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"прочитать позицию заказа %q: %w",
				orderUUID,
				err,
			)
		}

		items = append(items, converter.FromRecordToModelOrderItem(itemRecord))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"прочитать строки позиций заказа %q: %w",
			orderUUID,
			err,
		)
	}

	return items, nil
}
