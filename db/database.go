package db

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/lyneq/mailapi/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Username   string `json:"username" gorm:"uniqueIndex;not null"`
	Password   string `json:"-" gorm:"not null"`
	Role       string `json:"role" gorm:"not null;default:User"`
	IsVerified bool   `json:"is_verified" gorm:"default:false"`
}

func Init() {
	// Get database configuration from config.ini
	dbPath := config.GetDatabasePath()
	if dbPath == "" {
		dbPath = "./db/default.db" // Default path if not specified in config
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		_ = fmt.Errorf("failed to connect database: %v\n", err)
		log.Panicf("failed to connect database %v", err)
	}

	DB = db

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		_ = fmt.Errorf("failed to migrate database: %v\n", err)
		return
	}

	fmt.Println("Database connected")

	// Read
	//var product User
	//db.First(&product, 1)                  // find product with integer primary key
	//db.First(&product, "Role = ?", "User") // find product with code D42
}
