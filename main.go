package main

import (
	"github.com/dennisenwiya/Go-AUTH/controllers"
	"github.com/dennisenwiya/Go-AUTH/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}
func main() {
	r := gin.Default()
	r.POST("/user", controllers.UserRegister)
	r.POST("/login", controllers.UserLogin)

	r.Run() // listen and serve on 0.0.0.0:8080
}
