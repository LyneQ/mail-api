package db

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        string `json:"id" gorm:"primaryKey;not null"`
	Username  string `json:"username" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"`
	Role      string `json:"role" gorm:"not null" default:"User"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
}

func Init() {
	db, err := gorm.Open(sqlite.Open("./db/dev.db"), &gorm.Config{})
	if err != nil {
		_ = fmt.Errorf("failed to connect database: %v\n", err)
		log.Panicf("failed to connect database %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		_ = fmt.Errorf("failed to migrate database: %v\n", err)
		return
	}

	fmt.Println("Database connected")
	// db.Create(&User{Username: "Default User", Password: "123"})

	// Read
	//var product User
	//db.First(&product, 1)                  // find product with integer primary key
	//db.First(&product, "Role = ?", "User") // find product with code D42
}
