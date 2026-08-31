package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"medicalBooking/models"
	"medicalBooking/repository"
	"medicalBooking/views"
)

// ReservationController orquesta la creación/edición de reservas,
// validando reglas de negocio como el solapamiento de horarios.
type ReservationController struct {
	repo         repository.ReservationRepository
	userRepo     repository.UserRepository
	resourceRepo repository.ResourceRepository
}

func NewReservationController(
	repo repository.ReservationRepository,
	userRepo repository.UserRepository,
	resourceRepo repository.ResourceRepository,
) *ReservationController {
	return &ReservationController{
		repo:         repo,
		userRepo:     userRepo,
		resourceRepo: resourceRepo,
	}
}

// POST /api/reservations
func (c *ReservationController) Create(ctx *gin.Context) {
	var req views.CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !req.EndTime.After(req.StartTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "end_time debe ser posterior a start_time"})
		return
	}

	if _, err := c.userRepo.FindByID(req.UserID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	resource, err := c.resourceRepo.FindByID(req.ResourceID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "recurso no encontrado"})
		return
	}
	if !resource.Active {
		ctx.JSON(http.StatusConflict, gin.H{"error": "el recurso no está disponible actualmente"})
		return
	}

	overlap, err := c.repo.HasOverlap(req.ResourceID, req.StartTime, req.EndTime, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo validar disponibilidad"})
		return
	}
	if overlap {
		ctx.JSON(http.StatusConflict, gin.H{"error": "el recurso ya está reservado en ese horario"})
		return
	}

	reservation := models.Reservation{
		UserID:     req.UserID,
		ResourceID: req.ResourceID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Notes:      req.Notes,
		Status:     models.StatusPending,
	}

	if err := c.repo.Create(&reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear la reserva"})
		return
	}

	created, err := c.repo.FindByID(reservation.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reserva creada pero no se pudo recuperar"})
		return
	}

	ctx.JSON(http.StatusCreated, views.ToReservationResponse(*created))
}

// GET /api/reservations
func (c *ReservationController) FindAll(ctx *gin.Context) {
	reservations, err := c.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener las reservas"})
		return
	}
	ctx.JSON(http.StatusOK, views.ToReservationResponseList(reservations))
}

// GET /api/reservations/:id
func (c *ReservationController) FindByID(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	reservation, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, views.ToReservationResponse(*reservation))
}

// GET /api/users/:id/reservations
func (c *ReservationController) FindByUser(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	reservations, err := c.repo.FindByUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener las reservas del usuario"})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

// GET /api/resources/:id/reservations
func (c *ReservationController) FindByResource(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	reservations, err := c.repo.FindByResource(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener las reservas del recurso"})
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}

// PUT /api/reservations/:id
func (c *ReservationController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	reservation, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
		return
	}

	var req views.UpdateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStart := reservation.StartTime
	newEnd := reservation.EndTime
	if req.StartTime != nil {
		newStart = *req.StartTime
	}
	if req.EndTime != nil {
		newEnd = *req.EndTime
	}
	if !newEnd.After(newStart) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "end_time debe ser posterior a start_time"})
		return
	}

	if req.StartTime != nil || req.EndTime != nil {
		overlap, err := c.repo.HasOverlap(reservation.ResourceID, newStart, newEnd, reservation.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo validar disponibilidad"})
			return
		}
		if overlap {
			ctx.JSON(http.StatusConflict, gin.H{"error": "el recurso ya está reservado en ese horario"})
			return
		}
		reservation.StartTime = newStart
		reservation.EndTime = newEnd
	}

	if req.Status != "" {
		reservation.Status = req.Status
	}
	if req.Notes != "" {
		reservation.Notes = req.Notes
	}

	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar la reserva"})
		return
	}

	updated, _ := c.repo.FindByID(reservation.ID)
	ctx.JSON(http.StatusOK, views.ToReservationResponse(*updated))
}

// PATCH /api/reservations/:id/cancel
func (c *ReservationController) Cancel(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	reservation, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
		return
	}

	reservation.Status = models.StatusCancelled
	if err := c.repo.Update(reservation); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo cancelar la reserva"})
		return
	}

	ctx.JSON(http.StatusOK, views.ToReservationResponse(*reservation))
}

// DELETE /api/reservations/:id
func (c *ReservationController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	if _, err := c.repo.FindByID(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
		return
	}

	if err := c.repo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo eliminar la reserva"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
