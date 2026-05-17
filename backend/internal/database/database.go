package database

import (
	"log"

	"kare-rehber/backend/internal/config"
	"kare-rehber/backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	cfg := config.App

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Coach{},
		&models.Coordinator{},
		&models.CoachStudent{},
		&models.CoordinatorStudent{},
		&models.Week{},
		&models.Meeting{},
		&models.MeetingLog{},
		&models.Message{},
		&models.SMSLog{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Println("Database connected and migrated.")
	DB = db
}
