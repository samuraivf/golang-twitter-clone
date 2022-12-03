package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "gateway/docs"
)

type Handler struct {}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitServer() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.isUnauthorized, func(c *gin.Context) {
				h.signIn(c, h.createTokens)
			})
			auth.GET("/logout", h.logout)
			auth.GET("/refresh", func(c *gin.Context) {
				h.refresh(c, h.createTokens)
			})
		}

		user := api.Group("/user", h.isAuthorized)
		{
			user.PUT("/edit-profile", h.editProfile)
			user.POST("/add-image", h.addImage)
			user.GET("/:username", h.getUserByUsername)
			user.POST("/subscribe/:id", h.Subscribe)
			user.POST("/unsubscribe/:id", h.Unsubscribe)
			user.GET("/messages", h.getUserMessages)
		}

		tweet := api.Group("/tweet", h.isAuthorized)
		{
			tweet.POST("/create", h.createTweet)
			tweet.GET("/:id", h.getTweetById)
			tweet.GET("/user-tweets/:userId", h.getUserTweets)
			tweet.PUT("/update", h.updateTweet)
			tweet.DELETE("/:id", h.deleteTweet)
			tweet.GET("/like/:id", h.likeTweet)
			tweet.GET("/unlike/:id", h.unlikeTweet)
		}

		comment := api.Group("/comment", h.isAuthorized)
		{
			comment.POST("/create", h.createComment)
			comment.GET("/:id", h.getCommentById)
			comment.PUT("/update", h.updateComment)
			comment.DELETE("/:id", h.deleteComment)
			comment.GET("/like/:id", h.likeComment)
			comment.GET("/unlike/:id", h.unlikeComment)
		}

		tag := api.Group("/tag", h.isAuthorized)
		{
			tag.GET("/with-tweets/:id", h.getTagByIdWithTweets)
			tag.GET("/:name", h.getTagByName)
		}
	}

	return router
}
