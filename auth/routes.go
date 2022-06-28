package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func loginForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login", gin.H{})
}

func AddAuthRoutes(router *gin.Engine) {
	router.GET("/login", loginForm)
	// router.POST("/auth", auth)
}
