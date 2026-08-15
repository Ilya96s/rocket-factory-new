package order

type service struct {
	orderRepo       OrderRepository
	orderItemRepo   OrderItemRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	txManager       TxManager
}

func NewService(orderRepo OrderRepository, orderItemRepo OrderItemRepository, inventoryClient InventoryClient, paymentClient PaymentClient, manager TxManager) *service {
	return &service{
		orderRepo:       orderRepo,
		orderItemRepo:   orderItemRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		txManager:       manager,
	}
}
