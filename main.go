package main

import (
	"fmt"
	"github.com/lyneq/mailapi/api"
	"github.com/lyneq/mailapi/config"
	"github.com/lyneq/mailapi/db"
	"github.com/lyneq/mailapi/internal/session"
	"os"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	db.Init()
	session.Init(db.DB, false)
	api.Init()
}
