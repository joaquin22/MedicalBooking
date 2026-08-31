package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"medicalBooking/models"
	"medicalBooking/repository"
	"medicalBooking/views"
)

// ResourceController maneja lo que se reserva: médicos/consultorios, salas, etc.
type ResourceController struct {
	repo repository.ResourceRepository
}

func NewResourceController(repo repository.ResourceRepository) *ResourceController {
	return &ResourceController{repo: repo}
}

// POST /api/resources
func (c *ResourceController) Create(ctx *gin.Context) {
	var req views.CreateResourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	capacity := req.Capacity
	if capacity <= 0 {
		capacity = 1
	}

	resource := models.Resource{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Location:    req.Location,
		Capacity:    capacity,
		Active:      true,
	}

	if err := c.repo.Create(&resource); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el recurso"})
		return
	}

	ctx.JSON(http.StatusCreated, views.ToResourceResponse(resource))
}

// GET /api/resources
func (c *ResourceController) FindAll(ctx *gin.Context) {
	resources, err := c.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los recursos"})
		return
	}
	ctx.JSON(http.StatusOK, views.ToResourceResponseList(resources))
}

// GET /api/resources/:id
func (c *ResourceController) FindByID(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	resource, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "recurso no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, views.ToResourceResponse(*resource))
}

// PUT /api/resources/:id
func (c *ResourceController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	resource, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "recurso no encontrado"})
		return
	}

	var req views.UpdateResourceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		resource.Name = req.Name
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	if req.Location != "" {
		resource.Location = req.Location
	}
	if req.Capacity > 0 {
		resource.Capacity = req.Capacity
	}
	if req.Active != nil {
		resource.Active = *req.Active
	}

	if err := c.repo.Update(resource); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar el recurso"})
		return
	}

	ctx.JSON(http.StatusOK, views.ToResourceResponse(*resource))
}

// DELETE /api/resources/:id
func (c *ResourceController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	if _, err := c.repo.FindByID(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "recurso no encontrado"})
		return
	}

	if err := c.repo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo eliminar el recurso"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
