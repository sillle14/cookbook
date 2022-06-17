package main

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func formatObjId(objId primitive.ObjectID) string {
	return objId.Hex()
}

func formatInstructions(raw string) template.HTML {
	lb := regexp.MustCompile("[\r\n]+")
	steps := lb.Split(raw, -1)
	ret := ""
	for _, step := range steps {
		ret += "<li>" + step + "</li>\n"
	}
	return template.HTML(ret)
}

func CreateMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	funcMap := template.FuncMap{
		"formatObjId":        formatObjId,
		"formatInstructions": formatInstructions,
	}

	layouts, err := filepath.Glob("./templates/*.html.tmpl")
	if err != nil {
		panic(err.Error())
	}

	for _, layout := range layouts {
		layoutName := strings.Split(filepath.Base(layout), ".")[0]
		if layoutName != "base" {
			r.AddFromFilesFuncs(layoutName, funcMap, "templates/base.html.tmpl", layout)
		}
	}
	return r
}
