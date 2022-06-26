package setup

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/sillle14/soups-up/recipe"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var lb *regexp.Regexp = regexp.MustCompile(`[\r\n]+`)
var ingRegexp *regexp.Regexp = regexp.MustCompile(`(?P<first>\d+\s\w+)\s\+\s(?P<second>\d+\s\w+)(\sof)?\s(?P<ing>\w+)`)
var urlRegexp *regexp.Regexp = regexp.MustCompile(`https?://`)

func formatObjId(objId primitive.ObjectID) string {
	return objId.Hex()
}

func boldIngredients(recipe recipe.Recipe) string {
	ingredients := lb.Split(recipe.Ingredients, -1)
	instructions := recipe.Instructions
	for _, ingredient := range ingredients {
		if strings.Contains(ingredient, "+") {
			match := ingRegexp.FindStringSubmatch(ingredient)
			result := make(map[string]string)
			for i, name := range ingRegexp.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			for _, idx := range []string{"first", "second"} {
				parsedIng := result[idx] + " " + result["ing"]
				instructions = strings.ReplaceAll(instructions, parsedIng, "<b>"+parsedIng+"</b>")
			}
		} else {
			instructions = strings.ReplaceAll(instructions, ingredient, "<b>"+ingredient+"</b>")
		}
	}
	return instructions
}

func formatList(raw string) template.HTML {
	steps := lb.Split(raw, -1)
	ret := ""
	for _, step := range steps {
		ret += "<li>" + step + "</li>\n"
	}
	return template.HTML(ret)
}

func isLink(text string) bool {
	return urlRegexp.MatchString(text)
}

func CreateMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	funcMap := template.FuncMap{
		"formatObjId":     formatObjId,
		"formatList":      formatList,
		"boldIngredients": boldIngredients,
		"isLink":          isLink,
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
