package main

import (
	"github.com/lyneq/mailapi/api"
	"github.com/lyneq/mailapi/db"
)

func main() {
	db.Init()
	api.Init()
}
