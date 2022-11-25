package github

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Github struct {
	clientGQL *githubv4.Client
	client    *github.Client
}

func New(ctx context.Context, token string) (*Github, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	clientREST := github.NewClient(tc)
	clientGQL := githubv4.NewClient(tc)

	return &Github{
		clientGQL,
		clientREST,
	}, nil
}
