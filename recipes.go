// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package cookbook

import (
	"fmt"

	"github.com/reviewpad/cookbook/codehost"
	"github.com/reviewpad/cookbook/recipes"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
)

type NewRecipe func(handler.TargetEntity, codehost.Codehost, collector.Collector) (recipes.Recipe, error)

var recs map[string]NewRecipe = map[string]NewRecipe{
	"size": func(te handler.TargetEntity, ch codehost.Codehost, co collector.Collector) (recipes.Recipe, error) {
		return recipes.NewSizeRecipe(te, ch, co)
	},
}

func GetRecipeByName(name string, te handler.TargetEntity, ch codehost.Codehost, co collector.Collector) (recipes.Recipe, error) {
	init, ok := recs[name]
	if !ok {
		return nil, fmt.Errorf(`"%s" recipe not found`, name)
	}

	return init(te, ch, co)
}
