package main

import (
	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/mr-emerald-wolf/yantra-backend/models"
)

func init() {
	config, _ := initializers.LoadConfig(".")
	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
}
