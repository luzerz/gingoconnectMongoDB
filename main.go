package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luzerz/apijobtest/app/configs"
	"github.com/luzerz/apijobtest/app/routes"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()
	routes.UserRoute(router) //add this

	router.Run("localhost:80")
}
