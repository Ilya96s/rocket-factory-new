package record // TODO какое название пакета правильнее использовать converter или order ?

import "time"

// TODO изменить на тип указатель у тех типов где мб nil
type OrderRecord struct {
	UUID            string
	HulUUID         string
	EngineUUID      string
	ShieldUUID      string
	WeaponUUID      string
	TotalPrice      int64
	TransactionUUID string
	PaymentMethod   string
	Status          string
	CreatedAt       time.Time
}
