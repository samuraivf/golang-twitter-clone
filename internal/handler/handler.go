package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samuraivf/twitter-clone/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitServer() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.isUnauthorized, h.signIn)
			auth.GET("/logout", h.logout)
			auth.GET("/refresh", h.refresh)
		}

		user := api.Group("/user", h.isAuthorized)
		{
			user.PUT("/edit-profile", h.editProfile)
			user.POST("/add-image", h.addImage)
			user.GET("/:username", h.getUserByUsername)
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
	}

	return router
}

func (h *Handler) Hello(c *gin.Context) {
	tokeData, err := getUserData(c)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":       tokeData.UserId,
		"username": tokeData.Username,
	})
}
