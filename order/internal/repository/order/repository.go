package order

import (
	"sync"

	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/record"
	"github.com/Ilya96s/rocket-factory-new/order/internal/service/order"
)

var _ order.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[string]record.OrderRecord
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]record.OrderRecord),
	}
}
