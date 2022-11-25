// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes

import "context"

type Recipe interface {
	Run(context.Context) error
}
