package part

import (
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/service/part"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ part.PartRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
