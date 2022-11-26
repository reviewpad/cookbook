// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

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
