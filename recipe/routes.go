package recipe

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sillle14/soups-up/auth"
	"github.com/sillle14/soups-up/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func index(ctx *gin.Context) {
	opts := options.Find().SetProjection(bson.D{{Key: "name", Value: 1}})
	opts.SetLimit(100)
	opts.SetCollation(&options.Collation{Locale: "en"}) // enables case-insensitive sort
	opts.SetSort(bson.D{{Key: "name", Value: 1}})
	cursor, err := db.RecipesCollection.Find(ctx.Request.Context(), bson.D{}, opts)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var recipes []Recipe
	if err = cursor.All(ctx.Request.Context(), &recipes); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Ensure new recipes appear in the index.
	ctx.Header("Cache-Control", "no-cache, must-revalidate, no-store")
	ctx.HTML(http.StatusOK, "index", gin.H{
		"Recipes": recipes,
	})
}

func single(ctx *gin.Context) {
	id := ctx.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var recipe Recipe
	err = db.RecipesCollection.FindOne(ctx.Request.Context(), bson.M{"_id": objId}).Decode(&recipe)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "recipe", gin.H{
		"Recipe": recipe,
	})
}

func newForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form", gin.H{
		"Title": "New Recipe",
	})
}

func new(ctx *gin.Context) {
	var recipe RecipeForm
	err := ctx.ShouldBind(&recipe)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	res, err := db.RecipesCollection.InsertOne(ctx.Request.Context(), recipe)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/recipes/"+res.InsertedID.(primitive.ObjectID).Hex())
}

func editForm(ctx *gin.Context) {
	id := ctx.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var recipe Recipe
	err = db.RecipesCollection.FindOne(ctx.Request.Context(), bson.M{"_id": objId}).Decode(&recipe)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "form", gin.H{
		"Title":  "Edit Recipe",
		"Recipe": recipe,
		"Edit":   true,
	})
}

func edit(ctx *gin.Context) {
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
	_, err = db.RecipesCollection.UpdateByID(ctx.Request.Context(), objId, update)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/recipes/"+id)
}

func AddRecipeRoutes(router *gin.Engine) {
	recipes := router.Group("/recipes", auth.Middleware)
	recipes.GET("/", index)
	recipes.GET("/:id", single)
	recipes.GET("/new", newForm)
	recipes.POST("/", new)
	recipes.GET("/:id/edit", editForm)
	recipes.POST("/:id", edit)
}
