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
