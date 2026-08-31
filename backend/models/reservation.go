package models

import (
	"time"

	"gorm.io/gorm"
)

// ReservationStatus controla el ciclo de vida de una reserva
type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusConfirmed ReservationStatus = "confirmed"
	StatusCancelled ReservationStatus = "cancelled"
	StatusCompleted ReservationStatus = "completed"
)

// Reservation vincula a un User con un Resource en un rango de tiempo
type Reservation struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	UserID    uint              `gorm:"not null;index" json:"user_id"`
	ResourceID uint             `gorm:"not null;index" json:"resource_id"`
	StartTime time.Time         `gorm:"not null;index" json:"start_time"`
	EndTime   time.Time         `gorm:"not null" json:"end_time"`
	Status    ReservationStatus `gorm:"size:20;not null;default:pending" json:"status"`
	Notes     string            `gorm:"size:500" json:"notes"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`

	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Resource Resource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
}

func (Reservation) TableName() string {
	return "reservations"
}
