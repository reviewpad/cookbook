package github

import (
	"context"
	"errors"

	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/cookbook/codehost"
)

// CreateLabels batch creates labels ignoring any that already exist
func (gh *Github) CreateLabels(ctx context.Context, owner, repo string, labels []codehost.Label) error {
	for _, label := range labels {
		_, _, err := gh.client.Issues.CreateLabel(ctx, owner, repo, &github.Label{
			Name:        &label.Name,
			Color:       &label.Color,
			Description: &label.Description,
		})
		if err != nil {
			errorResponse := &github.ErrorResponse{}
			if errors.As(err, &errorResponse) {
				if len(errorResponse.Errors) > 0 && errorResponse.Errors[0].Code == "already_exists" {
					continue
				}
			}

			return err
		}
	}

	return nil
}

func (gh *Github) SetLabels(ctx context.Context, owner, repo string, number int, labels []string) error {
	_, _, err := gh.client.Issues.ReplaceLabelsForIssue(ctx, owner, repo, number, labels)
	return err
}
