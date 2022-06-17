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
// - better follow rest, all endpoints should be /recipes
//    - maybe / redirects to /recipes?
// - form and post endpoint
// - super simple auth, ask for a password and then send a signed JWT in a cookie. Check on all, else redirect to `/login`
// - figure out the log message about trust all proxies
// - Home button
// - content -> instructions
// - name -> something else

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
	r.AddFromFilesFuncs("form", funcMap, "templates/base.html.tmpl", "templates/form.html.tmpl")
	return r
}

func main() {

	// config
	ConnectDB()

	router := gin.Default()

	router.HTMLRender = createMyRender()

	router.Static("/assets", "./assets")
	router.GET("/", func(ctx *gin.Context) {
		opts := options.Find().SetProjection(bson.D{{Key: "name", Value: 1}})
		opts.SetLimit(100)
		cursor, err := RecipesCollection.Find(ctx.Request.Context(), bson.D{}, opts)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var recipes []Recipe
		if err = cursor.All(ctx.Request.Context(), &recipes); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// TODO: This means that new recipes show in the index when added.
		// TODO: How to just add this header after a new one is added?
		ctx.Header("Cache-Control", "no-cache, must-revalidate, no-store")
		ctx.HTML(http.StatusOK, "index", gin.H{
			"Recipes": recipes,
		})
	})

	router.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var recipe Recipe
		err = RecipesCollection.FindOne(ctx.Request.Context(), bson.M{"_id": objId}).Decode(&recipe)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.HTML(http.StatusOK, "recipe", gin.H{
			"Recipe": recipe,
		})
	})

	router.GET("/new", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "form", gin.H{
			"Title": "New Recipe",
		})
	})

	router.POST("/recipes", func(ctx *gin.Context) {
		var recipe RecipeForm
		err := ctx.ShouldBind(&recipe)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		res, err := RecipesCollection.InsertOne(ctx.Request.Context(), recipe)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Redirect(http.StatusFound, "/" + res.InsertedID.(primitive.ObjectID).Hex())
	})

	router.GET("/:id/edit", func(ctx *gin.Context) {
		id := ctx.Param("id")
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var recipe Recipe
		err = RecipesCollection.FindOne(ctx.Request.Context(), bson.M{"_id": objId}).Decode(&recipe)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.HTML(http.StatusOK, "form", gin.H{
			"Title": "Edit Recipe",
			"Recipe": recipe,
			"Edit": true,
		})
	})

	router.POST("/recipes/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var recipe RecipeForm
		err = ctx.ShouldBind(&recipe)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		update := bson.D{{Key: "$set", Value: recipe}}
		_, err = RecipesCollection.UpdateByID(ctx.Request.Context(), objId, update)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Redirect(http.StatusFound, "/" + id)
	})

	router.Run()
}
