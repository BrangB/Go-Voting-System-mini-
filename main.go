package main

import (
	"github.com/brangb/go_voting_system/config"
	"github.com/brangb/go_voting_system/controllers"
	"github.com/brangb/go_voting_system/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariable()
	config.ConnectToDB()
	config.SyncDatabases()
}

func main() {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/ping", controllers.Ping)
		apiV1.POST("/register", controllers.RegisterUser)
		apiV1.POST("/login", controllers.Login)
		apiV1.GET("/validate", middlewares.CheckAuth, controllers.Validate)
		apiV1.POST("/poll", middlewares.CheckAuth, controllers.CreatePoll)
	}

	r.Run()
}
