package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func getJWTKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return HMACSecret, nil
}

func Middleware(ctx *gin.Context) {
	cookie, err := ctx.Cookie("session")
        if err != nil {
		ctx.Redirect(http.StatusFound, "/login")
        }
	
	token, err := jwt.Parse(cookie, getJWTKey)
	if err == nil && token.Valid {
		ctx.Next()
	} else {
		ctx.Redirect(http.StatusFound, "/login")
	}
}
