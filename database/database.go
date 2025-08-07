package database

import (
	"MSAgent/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

// InitializeDB initializes the SQLite database and performs auto-migration.
func InitializeDB() *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open("monitoring.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.CPUInfo{}, &models.MemoryInfo{}, &models.DiskInfo{}, &models.NetworkInfo{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated.")

	return db
}
