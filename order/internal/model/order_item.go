package model

import "time"

// OrderItem - модель описания деталей заказа на уровне сервисного слоя
type OrderItem struct {
	UUID      string
	OrderUUID string
	PartUUID  string
	PartType  PartType
	Price     int64
	CreatedAt time.Time
}

// PartType Тип детали
type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)
