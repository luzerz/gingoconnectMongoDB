package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luzerz/apijobtest/app/controllers"
)

func UserRoute(router *gin.Engine) {
	router.POST("/internal-insert", controllers.CreateUser())
	router.GET("/internal-read", controllers.GetAUser())

}
