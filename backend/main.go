package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"medicalBooking/config"
	"medicalBooking/controllers"
	"medicalBooking/repository"
	"medicalBooking/routes"
)

func main() {
	// Carga variables de entorno desde .env (si no existe, sigue con las del sistema)
	if err := godotenv.Load(); err != nil {
		log.Println("no se encontró archivo .env, usando variables de entorno del sistema")
	}

	db := config.ConnectDatabase()

	// --- Inyección de dependencias ---
	// repository -> controller -> router
	userRepo := repository.NewUserRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	userController := controllers.NewUserController(userRepo)
	resourceController := controllers.NewResourceController(resourceRepo)
	reservationController := controllers.NewReservationController(reservationRepo, userRepo, resourceRepo)

	router := routes.SetupRouter(userController, resourceController, reservationController)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("servidor corriendo en http://localhost:%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("error al iniciar el servidor: %v", err)
	}
}
