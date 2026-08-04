package v1

import inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"

type handler struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService PartService
}

func New(partService PartService) *handler {
	return &handler{
		partService: partService,
	}
}
