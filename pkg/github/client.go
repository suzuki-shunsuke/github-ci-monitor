package github

import (
	"context"

	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

func New(ctx context.Context, token string) Client {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	return Client{
		client: github.NewClient(tc),
	}
}

func (client Client) GetCombinedStatus(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) (*github.CombinedStatus, *github.Response, error) {
	return client.client.Repositories.GetCombinedStatus(ctx, owner, repo, ref, opts)
}

func (client Client) ListCheckRunsForRef(ctx context.Context, owner, repo, ref string, opts *github.ListCheckRunsOptions) (*github.ListCheckRunsResults, *github.Response, error) {
	return client.client.Checks.ListCheckRunsForRef(ctx, owner, repo, ref, opts)
}

func (client Client) ListCheckSuitesForRef(ctx context.Context, owner, repo, ref string, opts *github.ListCheckSuiteOptions) (*github.ListCheckSuiteResults, *github.Response, error) {
	return client.client.Checks.ListCheckSuitesForRef(ctx, owner, repo, ref, opts)
}
