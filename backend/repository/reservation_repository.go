package repository

import (
	"time"

	"gorm.io/gorm"

	"medicalBooking/models"
)

type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	FindAll() ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	FindByUser(userID uint) ([]models.Reservation, error)
	FindByResource(resourceID uint) ([]models.Reservation, error)
	Update(reservation *models.Reservation) error
	Delete(id uint) error
	// HasOverlap verifica si ya existe una reserva activa que se cruce
	// en el tiempo para el mismo recurso (evita doble booking).
	HasOverlap(resourceID uint, start, end time.Time, excludeID uint) (bool, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Resource").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Preload("User").Preload("Resource").First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindByUser(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Resource").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByResource(resourceID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Where("resource_id = ?", resourceID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *reservationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Reservation{}, id).Error
}

// HasOverlap: dos rangos [start,end) se solapan si start < otro.end AND end > otro.start.
// Se excluyen reservas canceladas y, opcionalmente, la propia reserva (útil al editar).
func (r *reservationRepository) HasOverlap(resourceID uint, start, end time.Time, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Reservation{}).
		Where("resource_id = ?", resourceID).
		Where("status <> ?", models.StatusCancelled).
		Where("start_time < ? AND end_time > ?", end, start)

	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
