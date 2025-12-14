package repositories

import (
	"fmt"
	"log"
	"vasvault/internal/models"

	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Connect() (*DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL tidak diatur")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.File{}, &models.FileShare{}, &models.Category{}, &models.PublicLink{}, &models.Workspace{}, &models.WorkspaceMember{}); err != nil {
		log.Printf("Gagal melakukan migrasi: %v", err)
		return &DB{db}, err
	}

	return &DB{db}, nil
}
