package order_item

import (
	"context"
	"fmt"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/converter"
	"github.com/Masterminds/squirrel"
)

func (r *repository) Create(ctx context.Context, items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	queryBuilder := squirrel.Insert("order_items").
		Columns("uuid", "order_uuid", "part_uuid", "part_type", "price", "created_at").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range items {
		recordItem := converter.FromModelToRecordOrderItem(item)
		queryBuilder = queryBuilder.Values(
			recordItem.UUID,
			recordItem.OrderUUID,
			recordItem.PartUUID,
			recordItem.PartType,
			recordItem.Price,
			recordItem.CreatedAt,
		)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("сформировать запрос создания позиций заказа: %w", err)
	}

	if _, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("создать позиции заказа: %w", err)
	}

	return nil
}
