package part

import (
	"sync"

	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/record"
)

type repository struct {
	mu    sync.RWMutex
	parts map[string]record.PartRecord
}

func New() *repository {
	return &repository{
		parts: make(map[string]record.PartRecord),
	}
}
