package views

import (
	"time"

	"medicalBooking/models"
)

// CreateUserRequest es lo que el cliente envía para registrar un usuario
type CreateUserRequest struct {
	FullName string      `json:"full_name" binding:"required,min=3,max=150"`
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required,min=6"`
	Phone    string      `json:"phone"`
	Role     models.Role `json:"role"`
}

// UpdateUserRequest permite actualizar campos parciales de un usuario
type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"omitempty,min=3,max=150"`
	Phone    string `json:"phone"`
	Active   *bool  `json:"active"`
}

// UserResponse es lo que se expone al cliente (nunca incluye password_hash)
type UserResponse struct {
	ID        uint        `json:"id"`
	FullName  string      `json:"full_name"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
	Role      models.Role `json:"role"`
	Active    bool        `json:"active"`
	CreatedAt time.Time   `json:"created_at"`
}

// ToUserResponse mapea el modelo de dominio al DTO de salida
func ToUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
	}
}

func ToUserResponseList(users []models.User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, ToUserResponse(u))
	}
	return result
}
