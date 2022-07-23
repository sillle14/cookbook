package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sillle14/soups-up/auth"
	"github.com/sillle14/soups-up/db"
	"github.com/sillle14/soups-up/recipe"
	"github.com/sillle14/soups-up/setup"
)

// TODO:
// - Go home (from login) if logged in
// - use https://github.com/unrolled/secure to send to https all the time
// - add recipes and tweak bolding logic
// - search functionality
// - need css asset versioning to cache bust
// - make it pretty
// - tags

func main() {

	db.ConnectDB()
	router := setup.InitRouter()

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/recipes")
	})

	router.GET("/about", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "about", gin.H{})
	})

	recipe.AddRecipeRoutes(router)
	auth.AddAuthRoutes(router)

	router.Run()
}
