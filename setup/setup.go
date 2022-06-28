package setup

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.HTMLRender = CreateMyRender()
	router.Static("/assets", "./assets")
	return router
}
