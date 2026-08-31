package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"medicalBooking/controllers"

)

// SetupRouter receives the already constructed controllers (with their repositories injected)
// and defines all the API endpoints.
func SetupRouter(
	userController *controllers.UserController,
	resourceController *controllers.ResourceController,
	reservationController *controllers.ReservationController,
) *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware())

	api := router.Group("/api")
	{
		api.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "ok"})
		})

		users := api.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("", userController.FindAll)
			users.GET("/:id", userController.FindByID)
			users.PUT("/:id", userController.Update)
			users.DELETE("/:id", userController.Delete)
			users.GET("/:id/reservations", reservationController.FindByUser)
		}

		resources := api.Group("/resources")
		{
			resources.POST("", resourceController.Create)
			resources.GET("", resourceController.FindAll)
			resources.GET("/:id", resourceController.FindByID)
			resources.PUT("/:id", resourceController.Update)
			resources.DELETE("/:id", resourceController.Delete)
			resources.GET("/:id/reservations", reservationController.FindByResource)
		}

		reservations := api.Group("/reservations")
		{
			reservations.POST("", reservationController.Create)
			reservations.GET("", reservationController.FindAll)
			reservations.GET("/:id", reservationController.FindByID)
			reservations.PUT("/:id", reservationController.Update)
			reservations.PATCH("/:id/cancel", reservationController.Cancel)
			reservations.DELETE("/:id", reservationController.Delete)
		}
	}

	return router
}

// corsMiddleware agrega los encabezados CORS a las respuestas.
func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
