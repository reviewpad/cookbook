package cookbook

import (
	"github.com/reviewpad/cookbook/codehost"
	"github.com/reviewpad/cookbook/recipes"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
)

type NewRecipe func(handler.TargetEntity, codehost.Codehost, collector.Collector) (recipes.Recipe, error)

var Recipes map[string]NewRecipe = map[string]NewRecipe{
	"size": func(te handler.TargetEntity, ch codehost.Codehost, co collector.Collector) (recipes.Recipe, error) {
		return recipes.NewSizeRecipe(te, ch, co)
	},
}
