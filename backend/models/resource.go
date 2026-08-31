package models

import (
	"time"

	"gorm.io/gorm"
)

// ResourceType permite reutilizar el mismo modelo para distintos rubros
// (citas médicas, salas de reunión, canchas, equipos, etc.)
type ResourceType string

const (
	ResourceTypeDoctor ResourceType = "doctor"
	ResourceTypeRoom    ResourceType = "room"
	ResourceTypeOther   ResourceType = "other"
)

// Resource es "lo que se reserva": un médico/consultorio, una sala, etc.
type Resource struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Type        ResourceType   `gorm:"size:20;not null" json:"type"`
	Description string         `gorm:"size:255" json:"description"`
	Location    string         `gorm:"size:150" json:"location"`
	Capacity    int            `gorm:"default:1" json:"capacity"`
	Active      bool           `gorm:"default:true" json:"active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Reservations []Reservation `gorm:"foreignKey:ResourceID" json:"reservations,omitempty"`
}

func (Resource) TableName() string {
	return "resources"
}
