package order

type service struct {
	orderRepo       OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

func NewService(orderRepo OrderRepository, inventoryClient InventoryClient, paymentClient PaymentClient) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
