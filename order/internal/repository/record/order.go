package record

import "time"

type Order struct {
	UUID            string     `db:"uuid"`
	TotalPrice      int64      `db:"total_price"`
	Status          string     `db:"status"`
	TransactionUUID *string    `db:"transaction_uuid"`
	PaymentMethod   *string    `db:"payment_method"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}
