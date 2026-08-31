package views

import (
	"time"

	"medicalBooking/models"
)

type CreateReservationRequest struct {
	UserID     uint      `json:"user_id" binding:"required"`
	ResourceID uint      `json:"resource_id" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
	Notes      string    `json:"notes"`
}

type UpdateReservationRequest struct {
	StartTime *time.Time               `json:"start_time"`
	EndTime   *time.Time               `json:"end_time"`
	Status    models.ReservationStatus `json:"status"`
	Notes     string                   `json:"notes"`
}

type ReservationResponse struct {
	ID         uint                     `json:"id"`
	Status     models.ReservationStatus `json:"status"`
	StartTime  time.Time                `json:"start_time"`
	EndTime    time.Time                `json:"end_time"`
	Notes      string                   `json:"notes"`
	User       UserResponse             `json:"user"`
	Resource   ResourceResponse         `json:"resource"`
	CreatedAt  time.Time                `json:"created_at"`
}

func ToReservationResponse(res models.Reservation) ReservationResponse {
	return ReservationResponse{
		ID:        res.ID,
		Status:    res.Status,
		StartTime: res.StartTime,
		EndTime:   res.EndTime,
		Notes:     res.Notes,
		User:      ToUserResponse(res.User),
		Resource:  ToResourceResponse(res.Resource),
		CreatedAt: res.CreatedAt,
	}
}

func ToReservationResponseList(list []models.Reservation) []ReservationResponse {
	result := make([]ReservationResponse, 0, len(list))
	for _, res := range list {
		result = append(result, ToReservationResponse(res))
	}
	return result
}
