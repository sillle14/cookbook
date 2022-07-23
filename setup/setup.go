package setup

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func redirectToHTTPS() gin.HandlerFunc {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
	})
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()
}

func InitRouter() *gin.Engine {
	router := gin.Default()
	gin_mode, ok := os.LookupEnv("GIN_MODE")
	if ok && gin_mode == "release" {
		router.Use(redirectToHTTPS())
	}
	router.HTMLRender = CreateMyRender()
	router.Static("/assets", "./assets")
	router.SetTrustedProxies(nil)
	return router
}
