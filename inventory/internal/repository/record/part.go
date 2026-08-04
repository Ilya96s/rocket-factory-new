package record

import "time"

type PartRecord struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      string
	StockQuantity int64
	CreatedAt     time.Time
}
