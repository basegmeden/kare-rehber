package database

import (
	"log"

	"kare-rehber/backend/internal/config"
	"kare-rehber/backend/internal/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
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

	ensureAdmin(db)

	log.Println("Database connected and migrated.")
	DB = db
}

func ensureAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ensureAdmin: bcrypt error: %v", err)
		return
	}
	admin := models.User{
		Name:         "Admin",
		Surname:      "KARE",
		Phone:        "00000000000",
		City:         "",
		Role:         models.RoleAdmin,
		Username:     "admin",
		PasswordHash: string(hash),
		Status:       models.UserStatusActive,
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Printf("ensureAdmin: could not create admin user: %v", err)
		return
	}
	log.Println("Admin user created (admin / admin123)")
}
