package repository

import (
	"gorm.io/gorm"

	"medicalBooking/models"
)

type ResourceRepository interface {
	Create(resource *models.Resource) error
	FindAll() ([]models.Resource, error)
	FindByID(id uint) (*models.Resource, error)
	Update(resource *models.Resource) error
	Delete(id uint) error
}

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &resourceRepository{db: db}
}

func (r *resourceRepository) Create(resource *models.Resource) error {
	return r.db.Create(resource).Error
}

func (r *resourceRepository) FindAll() ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Find(&resources).Error
	return resources, err
}

func (r *resourceRepository) FindByID(id uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepository) Update(resource *models.Resource) error {
	return r.db.Save(resource).Error
}

func (r *resourceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Resource{}, id).Error
}
