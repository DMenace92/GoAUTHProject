package main

import (
	"github.com/dennisenwiya/Go-AUTH/initializers"

	"github.com/dennisenwiya/Go-AUTH/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}
func main() {
	initializers.DB.AutoMigrate(&models.User{})
}
