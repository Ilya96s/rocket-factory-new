package model

import "time"

type Part struct {
	UUID          string
	Name          string
	PartType      PartType
	Price         int64
	StockQuantity int64
	CreatedAt     time.Time
}
