// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes_test

import (
	"context"
	"testing"

	"github.com/reviewpad/cookbook"
	"github.com/reviewpad/cookbook/mocks"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/stretchr/testify/assert"
)

func TestSize(t *testing.T) {
	collector := mocks.NewCollector(t)
	codehost := mocks.NewCodehost(t)
	ctx := context.Background()

	tests := map[string]struct {
		targetEntity handler.TargetEntity
		wantError    error
		mock         func()
	}{}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.mock != nil {
				test.mock()
			}

			size, err := cookbook.GetRecipeByName("size", test.targetEntity, codehost, collector)
			assert.Nil(t, err)

			err = size.Run(ctx)
			assert.Equal(t, test.wantError, err)
		})
	}
}
