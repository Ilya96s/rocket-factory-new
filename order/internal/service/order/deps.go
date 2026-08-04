package order

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

// OrderRepository определяет контракт для работы с хранилищем заказов.
type OrderRepository interface {
	// Create - создает новый заказ или полностью заменяет существующий.
	Create(ctx context.Context, order model.Order) error

	// Get - возвращает заказ по UUID
	Get(ctx context.Context, uuid string) (model.Order, error)

	// Update - выполняет обновление/изменение заказа
	Update(ctx context.Context, order model.Order) error
}

// InventoryClient определяет контракты для работы с InventoryService
type InventoryClient interface {
	// ListParts - возвращает список деталей с возможностью фильтрации по типу или по конкретным UUID
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}

// PaymentClient определяет контракты для работы с PaymentService
type PaymentClient interface {
	// PayOrder - обрабатывает команду на оплату и возвращает transaction_uuid
	PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error)
}
