package main

import (
	"medicalBooking/config"
	"medicalBooking/models"
)

func main() {
	db := config.PostgresConnection()

	if err := db.AutoMigrate(&models.UserModel{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}
