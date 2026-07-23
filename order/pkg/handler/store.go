package handler

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// OrderStatus — внутренний статус заказа.
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — внутренний способ оплаты.
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order — внутренняя модель заказа.
type Order struct {
	OrderUUID  uuid.UUID
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID

	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID

	TotalPrice int64

	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod

	Status    OrderStatus
	CreatedAt time.Time
}

// orderStore — потокобезопасное in-memory-хранилище.
type orderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт пустое хранилище заказов.
func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// Get возвращает заказ по UUID.
func (s *orderStore) Get(orderUUID uuid.UUID) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	return order, ok
}

// Save сохраняет новый заказ или полностью заменяет существующий.
func (s *orderStore) Save(order Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}

// Update выполняет атомарное изменение заказа.
//
// Функция update вызывается под блокировкой.
func (s *orderStore) Update(
	orderUUID uuid.UUID,
	update func(order *Order) error,
) (Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return Order{}, false, nil
	}

	if err := update(&order); err != nil {
		return Order{}, true, err
	}

	s.orders[orderUUID] = order

	return order, true, nil
}
