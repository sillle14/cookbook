package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const MONTH = 2592000

type Password struct {
	Password string `form:"password"`
}

func loginForm(ctx *gin.Context) {
	_, failed := ctx.GetQuery("failed")
	ctx.HTML(http.StatusOK, "login", gin.H{"Failed": failed})
}

func login(ctx *gin.Context) {
	var password Password
	err := ctx.ShouldBind(&password)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if password.Password != ExpectedPassword {
		ctx.Redirect(http.StatusFound, "/login?failed=true")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().AddDate(0, 0, 30).Unix(),
	})
	tokenString, _ := token.SignedString(HMACSecret)
	ctx.SetCookie("session", tokenString, MONTH, "", "", true, true)
	ctx.Redirect(http.StatusFound, "/recipes/")
}

func AddAuthRoutes(router *gin.Engine) {
	router.GET("/login", loginForm)
	router.POST("/login", login)
}
