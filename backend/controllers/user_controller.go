package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"medicalBooking/models"
	"medicalBooking/repository"
	"medicalBooking/views"
)

// UserController maneja las peticiones HTTP y delega el acceso a datos al repository.
// No contiene lógica de GORM: eso vive exclusivamente en repository.
type UserController struct {
	repo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{repo: repo}
}

// POST /api/users
func (c *UserController) Create(ctx *gin.Context) {
	var req views.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := c.repo.FindByEmail(req.Email); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "ya existe un usuario con ese email"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar la contraseña"})
		return
	}

	role := req.Role
	if role == "" {
		role = models.RoleCustomer
	}

	user := models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hash),
		Phone:        req.Phone,
		Role:         role,
		Active:       true,
	}

	if err := c.repo.Create(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el usuario"})
		return
	}

	ctx.JSON(http.StatusCreated, views.ToUserResponse(user))
}

// GET /api/users
func (c *UserController) FindAll(ctx *gin.Context) {
	users, err := c.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los usuarios"})
		return
	}
	ctx.JSON(http.StatusOK, views.ToUserResponseList(users))
}

// GET /api/users/:id
func (c *UserController) FindByID(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	user, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, views.ToUserResponse(*user))
}

// PUT /api/users/:id
func (c *UserController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	user, err := c.repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	var req views.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Active != nil {
		user.Active = *req.Active
	}

	if err := c.repo.Update(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar el usuario"})
		return
	}

	ctx.JSON(http.StatusOK, views.ToUserResponse(*user))
}

// DELETE /api/users/:id
func (c *UserController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	if _, err := c.repo.FindByID(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	if err := c.repo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo eliminar el usuario"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// parseID es un helper compartido para leer y validar el parámetro :id de la URL
func parseID(ctx *gin.Context) (uint, error) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return 0, err
	}
	return uint(id), nil
}
