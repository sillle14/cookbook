package main

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO:
// - re organize
// - get css started
// - form and post endpoint

func formatObjId(objId primitive.ObjectID) string {
	return objId.Hex()
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	funcMap := template.FuncMap{
		"formatObjId": formatObjId,
	}
	// TODO: Loop through as there are more templates?
	r.AddFromFilesFuncs("index", funcMap, "templates/base.html.tmpl", "templates/index.html.tmpl")
	r.AddFromFilesFuncs("recipe", funcMap, "templates/base.html.tmpl", "templates/recipe.html.tmpl")
	return r
}

func main() {

	// config
	ConnectDB()

	router := gin.Default()

	router.HTMLRender = createMyRender()

	// router.Static("/assets", "./assets") Serve static files
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		opts := options.Find().SetProjection(bson.D{primitive.E{Key: "name", Value: "1"}})
		opts.SetLimit(100)
		cursor, err := RecipesCollection.Find(c.Request.Context(), bson.D{}, opts)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var recipes []RecipeLimited
		// TODO: Need to decode ID into string
		if err = cursor.All(c.Request.Context(), &recipes); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "index", gin.H{
			"Title":   "Home",
			"Recipes": recipes,
		})
	})

	router.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var recipe RecipeLimited
		err = RecipesCollection.FindOne(c.Request.Context(), bson.M{"_id": objId}).Decode(&recipe)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "recipe", gin.H{
			"Recipe": recipe,
		})
	})

	router.Run()
}
