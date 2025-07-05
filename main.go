package main

import (
	"github.com/lyneq/mailapi/api"
	"github.com/lyneq/mailapi/db"
	"github.com/lyneq/mailapi/internal/session"
)

func main() {
	db.Init()
	session.Init(db.DB, false)
	api.Init()
}
