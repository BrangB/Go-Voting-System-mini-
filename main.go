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

		//user --  CRUD
		apiV1.POST("/register", controllers.RegisterUser)
		apiV1.POST("/login", controllers.Login)
		apiV1.GET("/logout", controllers.Logout)
		apiV1.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)
		apiV1.PUT("/user/profile", middlewares.CheckAuth, controllers.UpdateProfile)

		//user -- validate
		apiV1.GET("/validate", middlewares.CheckAuth, controllers.Validate)

		//poll
		apiV1.GET("/polls", middlewares.CheckAuth, controllers.GetAllPolls)
		apiV1.GET("/poll/:id", middlewares.CheckAuth, controllers.GetPollById)
		apiV1.GET("/poll/:id/result", middlewares.CheckAuth, controllers.GetVoteResult)
		apiV1.POST("/poll", middlewares.CheckAuth, controllers.CreatePoll)
		apiV1.PUT("/poll/:id", middlewares.CheckAuth, controllers.UpdatePollByID)
		apiV1.DELETE("/poll/:id", middlewares.CheckAuth, controllers.DeletePollByID)

		//vote
		apiV1.GET("/vote/:id", middlewares.CheckAuth, controllers.VoteOption)
		apiV1.GET("/votes/option/:option_id", middlewares.CheckAuth, controllers.GetAllVotesByOptionId)
		apiV1.DELETE("/vote/:vote_id", middlewares.CheckAuth, controllers.RemoveVote)

	}

	r.Run()
}
