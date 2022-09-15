package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luzerz/apijobtest/app/controllers"
)

func UserRoute(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Available",
		})
	})
	router.POST("/internal-insert", controllers.CreateUser())
	router.GET("/internal-read", controllers.GetAUser())

}
