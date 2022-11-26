package github

import (
	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
)

type Github struct {
	clientGQL *githubv4.Client
	client    *github.Client
}

func New(clientREST *github.Client, clientGQL *githubv4.Client) *Github {
	return &Github{
		clientGQL,
		clientREST,
	}
}
