package main

import (
	"fmt"
	"log"

	"github.com/godongmin/open_talk/backend/internal/config"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/router"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := autoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	r := router.Setup(db, cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func setupDatabase(cfg *config.Config) (*gorm.DB, error) {
	if cfg.IsSQLite() {
		return gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{})
	}

	// For PostgreSQL support in production
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	// return gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{})
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.ChatRoom{},
		&model.ChatRoomMember{},
		&model.Message{},
		&model.Friend{},
	)
}
