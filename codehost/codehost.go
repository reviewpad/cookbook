// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package codehost

import (
	"context"
)

type Codehost interface {
	GetPRSizeData(context.Context, string, string, int) (*PRSizeData, error)
	CreateLabels(context.Context, string, string, []Label) error
	SetLabels(context.Context, string, string, int, []string) error
}

type Label struct {
	Name        string
	Color       string
	Description string
}

type PRSizeData struct {
	Changes uint64
	Labels  []string
}
