package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"medicalBooking/models"
)

// DB es la instancia global de la conexión a la base de datos.
// Se inyecta luego en repositorios (evita usar un singleton dentro de cada repo).
var DB *gorm.DB

// ConnectDatabase arma el DSN a partir de variables de entorno,
// abre la conexión con GORM y corre las migraciones automáticas.
func ConnectDatabase() *gorm.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "reservas_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("error al conectar a la base de datos: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("error al obtener sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := runMigrations(db); err != nil {
		log.Fatalf("error al ejecutar migraciones: %v", err)
	}

	DB = db
	log.Println("Conexión a PostgreSQL establecida y migraciones aplicadas")
	return db
}

// runMigrations centraliza el AutoMigrate de todos los modelos del sistema.
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Resource{},
		&models.Reservation{},
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
