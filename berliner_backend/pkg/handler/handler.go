package handler

import (
	"io"
	"os"

	"github.com/I1Asyl/berliner_backend/models"
	"github.com/I1Asyl/berliner_backend/pkg/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Handler struct which contains all services that are needed for the application
type Handler struct {
	services *services.Services
}

// NewHandler creates new Handler instance
func NewHandler(services *services.Services) *Handler {
	return &Handler{services: services}
}

// main page handler for user
func mainPage(ctx *gin.Context) {
	res, _ := ctx.Get("user")
	user := res.(models.User)

	ctx.JSON(200, gin.H{
		"username":  user.Username,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"id":        user.Id,
	})

}

// COR settings for router
func corSettings(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		// TODO: add domain
		return origin == "http://localhost:5173"
	}
	// possible methods
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}

	// make possible to share credientials
	config.AllowCredentials = true

	// allowed headers
	config.AllowHeaders = []string{"Origin", "Authorization"}

	router.Use(cors.New(config))
}

// InitRouter initializes router
func (h *Handler) InitRouter() *gin.Engine {
	// setting up logger
	//***
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	//***

	// creating a new router Engine
	router := gin.New()

	// setting up COR settings
	corSettings(router)

	// setting up middlewares
	router.Use(h.Logger())
	router.Use(gin.Recovery())

	// setting up authorization routes
	auth := router.Group("")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/login", h.login)
	}

	// setting up private routes
	private := router.Group("")
	{
		private.Use(h.AuthMiddleware())
		private.GET("", mainPage)

		private.GET("/channels", h.getChannels)
		private.POST("/channels", h.createChannel)
		private.PATCH("/channels", h.updateChannel)
		private.DELETE("/channels", h.deleteChannel)

		// post
		private.POST("/post", h.createPost)
		private.GET("/post", h.getPosts)
		private.DELETE("/post", h.deletePost)

		private.GET("/myPost", h.getMyChannelPosts)

		private.POST("/follow", h.follow)
		private.DELETE("/follow", h.unfollow)

		private.GET("/newPost", h.getNewPosts)

		private.GET("/following", h.getFollowing)

	}

	return router
}
