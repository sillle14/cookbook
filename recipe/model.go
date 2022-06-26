package recipe

import "go.mongodb.org/mongo-driver/bson/primitive"

// TODO: Add validation?
type RecipeForm struct {
	Name         string `json:"name" form:"name"`
	Source       string `json:"source" form:"source"`
	Ingredients  string `json:"ingredients" form:"ingredients"`
	Instructions string `json:"instructions" form:"instructions"`
}

type Recipe struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name" form:"name"`
	Source       string             `json:"source" form:"source"`
	Ingredients  string             `json:"ingredients" form:"ingredients"`
	Instructions string             `json:"instructions" form:"instructions"`
}
