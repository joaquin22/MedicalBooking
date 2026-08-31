package models

import (
	"time"

	"gorm.io/gorm"
)

// Role define el tipo de usuario dentro del sistema de reservas
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleStaff    Role = "staff"    // médico, encargado de sala, etc.
	RoleCustomer Role = "customer" // paciente, cliente que reserva
)

// User representa a cualquier persona que interactúa con el sistema
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	FullName     string         `gorm:"size:150;not null" json:"full_name"`
	Email        string         `gorm:"size:150;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Phone        string         `gorm:"size:30" json:"phone"`
	Role         Role           `gorm:"size:20;not null;default:customer" json:"role"`
	Active       bool           `gorm:"default:true" json:"active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Reservations []Reservation `gorm:"foreignKey:UserID" json:"reservations,omitempty"`
}

func (User) TableName() string {
	return "users"
}
