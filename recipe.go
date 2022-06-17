package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// TODO: This could maybe work with a full model using omitempty

type RecipeForm struct {
	Name string `json:"name" form:"name"`
	Content string `json:"content" form:"content"`
}

type Recipe struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" form:"name"`
	Content string          `json:"content" form:"content"`
}
