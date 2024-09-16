package controller

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v65/github"
	"github.com/suzuki-shunsuke/github-ci-monitor/pkg/datadog"
)

type Controller struct {
	GitHub  GitHub
	DataDog DataDog
	Stdout  io.Writer
	Stderr  io.Writer
}

type DataDog interface {
	Check(ctx context.Context, params datadog.ParamCheck) (*http.Response, error)
}

type GitHub interface {
	GetCombinedStatus(ctx context.Context, owner, repo, ref string, opts *github.ListOptions) (*github.CombinedStatus, *github.Response, error)
	ListCheckRunsForRef(ctx context.Context, owner, repo, ref string, opts *github.ListCheckRunsOptions) (*github.ListCheckRunsResults, *github.Response, error)
	ListCheckSuitesForRef(ctx context.Context, owner, repo, ref string, opts *github.ListCheckSuiteOptions) (*github.ListCheckSuiteResults, *github.Response, error)
}
