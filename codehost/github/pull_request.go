package github

import (
	"context"

	"github.com/reviewpad/cookbook/codehost"
	"github.com/shurcooL/githubv4"
)

type PRSizeDataQuery struct {
	Repository struct {
		PullRequest struct {
			Additions uint64 `graphql:"additions"`
			Deletions uint64 `graphql:"deletions"`
			Labels    struct {
				PageInfo struct {
					HasNextPage bool   `graphql:"hasNextPage"`
					EndCursor   string `graphql:"endCursor"`
				}
				Nodes []struct {
					Name string
				} `graphql:"nodes"`
			} `graphql:"labels(first: 100, after: $labelsEndCursor)"`
		} `graphql:"pullRequest(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

func (gh *Github) GetPRSizeData(ctx context.Context, owner, repo string, number int) (*codehost.PRSizeData, error) {
	var result PRSizeDataQuery
	labels := make([]string, 0)
	hasNextPage := true
	variables := map[string]interface{}{
		"owner":           githubv4.String(owner),
		"repo":            githubv4.String(repo),
		"number":          githubv4.Int(number),
		"labelsEndCursor": (*githubv4.String)(nil),
	}

	for hasNextPage {
		if err := gh.clientGQL.Query(ctx, &result, variables); err != nil {
			return nil, err
		}

		for _, node := range result.Repository.PullRequest.Labels.Nodes {
			labels = append(labels, node.Name)
		}

		variables["labelsEndCursor"] = result.Repository.PullRequest.Labels.PageInfo.EndCursor
		hasNextPage = result.Repository.PullRequest.Labels.PageInfo.HasNextPage
	}

	return &codehost.PRSizeData{
		Changes: (result.Repository.PullRequest.Additions + result.Repository.PullRequest.Deletions),
		Labels:  labels,
	}, nil
}
