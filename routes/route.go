package routes

import (
	"net/http"
	controller "notify/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/users", controller.GetAllTasks())
	router.POST("/user", controller.CreateTask())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Connected successfully to the server"})
	})
}
