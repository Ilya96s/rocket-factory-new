package v1

import (
	"context"

	"github.com/Ilya96s/rocket-factory-new/order/internal/model"
)

// OrderService - интерфейс сервисного слоя
type OrderService interface {
	// Create - создает новый заказ на основе выбранных компонентов корабля
	Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error)

	// Get - возвращает информацию о заказе по его UUID
	Get(ctx context.Context, uuid string) (model.Order, error)

	// Pay - проводит оплату ранее созданного заказа
	Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error)

	// Cancel - отменяет заказ, который еще не был оплачен
	Cancel(ctx context.Context, uuid string) error
}
