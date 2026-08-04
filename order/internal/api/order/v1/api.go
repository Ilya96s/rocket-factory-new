package v1

import orderv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/openapi/order/v1"

var _ orderv1.Handler = (*handler)(nil)

// handler - реализует HTTP API заказов
type handler struct {
	orderService OrderService
}

func New(orderService OrderService) *handler {
	return &handler{orderService: orderService}
}

func SetupServer(h *handler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}
