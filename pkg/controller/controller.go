package controller

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/github-ci-monitor/pkg/datadog"
	gh "github.com/suzuki-shunsuke/github-ci-monitor/pkg/github"
)

type Params struct {
	Repos         []Repo
	GitHubToken   string `yaml:"github_token"`
	DataDogAPIKey string `yaml:"datadog_api_key"`
	LogLevel      string `yaml:"log_level"`
	CheckName     string `yaml:"check_name"`
	Tags          map[string]string
}

type Repo struct {
	Owner      string
	Repo       string
	Ref        string
	Status     bool
	CheckRun   bool `yaml:"check_run"`
	CheckSuite bool `yaml:"check_suite"`
	Tags       map[string]string
}

const (
	kwFailure = "failure"
	kwTimeout = "timed_out"
)

func New(ctx context.Context, params Params) (Controller, Params, error) {
	if params.LogLevel != "" {
		lvl, err := logrus.ParseLevel(params.LogLevel)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"log_level": params.LogLevel,
			}).WithError(err).Error("the log level is invalid")
		}
		logrus.SetLevel(lvl)
	}

	if params.GitHubToken == "" {
		params.GitHubToken = os.Getenv("GITHUB_TOKEN")
		if params.GitHubToken == "" {
			params.GitHubToken = os.Getenv("GITHUB_ACCESS_TOKEN")
		}
	}
	if params.DataDogAPIKey == "" {
		params.DataDogAPIKey = os.Getenv("DATADOG_API_KEY")
	}

	return Controller{
		GitHub:  gh.New(ctx, params.GitHubToken),
		DataDog: datadog.New(params.DataDogAPIKey),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}, params, nil
}

var (
	errGitHubTokenRequired   = errors.New("GitHub Access Token is required")
	errDataDogAPIKeyRequired = errors.New("DataDog API Key is required")
	errOwnerRequired         = errors.New("owner is required")
	errRepoRequired          = errors.New("repo is required")
	errRefRequired           = errors.New("ref is required")
)

func (ctrl Controller) validateParams(params Params) error {
	if params.GitHubToken == "" {
		return errGitHubTokenRequired
	}
	if params.DataDogAPIKey == "" {
		return errDataDogAPIKeyRequired
	}
	for _, repo := range params.Repos {
		if repo.Owner == "" {
			return errOwnerRequired
		}
		if repo.Repo == "" {
			return errRepoRequired
		}
		if repo.Ref == "" {
			return errRefRequired
		}
	}
	return nil
}

func (ctrl Controller) CheckStatus(ctx context.Context, repo Repo) (bool, error) {
	status, _, err := ctrl.GitHub.GetCombinedStatus(ctx, repo.Owner, repo.Repo, repo.Ref, nil)
	if err != nil {
		return false, fmt.Errorf("get a combined status: %w", err)
	}
	return status.GetState() != kwFailure, nil
}

func (ctrl Controller) CheckRun(ctx context.Context, repo Repo) (bool, error) {
	checkRunResult, _, err := ctrl.GitHub.ListCheckRunsForRef(ctx, repo.Owner, repo.Repo, repo.Ref, nil)
	if err != nil {
		return false, fmt.Errorf("list check runs for a ref: %w", err)
	}
	for _, checkRun := range checkRunResult.CheckRuns {
		switch checkRun.GetConclusion() {
		case kwFailure, kwTimeout:
			return false, nil
		}
	}
	return true, nil
}

func (ctrl Controller) CheckSuite(ctx context.Context, repo Repo) (bool, error) {
	checkSuiteResult, _, err := ctrl.GitHub.ListCheckSuitesForRef(ctx, repo.Owner, repo.Repo, repo.Ref, nil)
	if err != nil {
		return false, fmt.Errorf("list check suites for a ref: %w", err)
	}
	for _, checkSuite := range checkSuiteResult.CheckSuites {
		switch checkSuite.GetConclusion() {
		case kwFailure, kwTimeout:
			return false, nil
		}
	}
	return true, nil
}

func (ctrl Controller) CheckRepo(ctx context.Context, repo Repo) (bool, error) {
	logE := logrus.WithFields(logrus.Fields{
		"repo":  repo.Repo,
		"owner": repo.Owner,
		"ref":   repo.Ref,
	})
	if repo.Status {
		logE.Info("check status")
		status, err := ctrl.CheckStatus(ctx, repo)
		if err != nil {
			return false, err
		}
		if !status {
			logE.Error("status is failure")
			return false, nil
		}
	}

	if repo.CheckRun {
		logE.Info("check run")
		run, err := ctrl.CheckRun(ctx, repo)
		if err != nil {
			return false, err
		}
		if !run {
			logE.Error("check_run is failure")
			return false, nil
		}
	}

	if repo.CheckSuite {
		logE.Info("check suite")
		suite, err := ctrl.CheckSuite(ctx, repo)
		if err != nil {
			return false, err
		}
		if !suite {
			logE.Error("check_suite is failure")
			return false, nil
		}
	}
	return true, nil
}

func (ctrl Controller) createTags(tagMaps ...map[string]string) []string {
	m := map[string]string{}
	for _, tagMap := range tagMaps {
		for k, v := range tagMap {
			m[k] = v
		}
	}
	tags := make([]string, 0, len(m))
	for k, v := range m {
		tags = append(tags, k+":"+v)
	}
	return tags
}

func (ctrl Controller) Run(ctx context.Context, params Params) error {
	if params.CheckName == "" {
		params.CheckName = "ci.ok"
	}
	if err := ctrl.validateParams(params); err != nil {
		return fmt.Errorf("parameters are invalid: %w", err)
	}

	for _, repo := range params.Repos {
		logE := logrus.WithFields(logrus.Fields{
			"repo":  repo.Repo,
			"owner": repo.Owner,
			"ref":   repo.Ref,
		})
		logE.Info("check a repo")
		b, err := ctrl.CheckRepo(ctx, repo)
		status := 0
		if err != nil {
			logE.WithError(err).Error("check a repo")
			status = 3
		} else if !b {
			status = 1
		}
		logE = logE.WithField("status", status)
		logE.Info("send a check to DataDog")

		_, err = ctrl.DataDog.Check(ctx, datadog.ParamCheck{ //nolint:bodyclose
			Status: status,
			Check:  params.CheckName,
			Tags: ctrl.createTags(
				params.Tags,
				repo.Tags,
				map[string]string{
					"owner": repo.Owner,
					"repo":  repo.Repo,
					"ref":   repo.Ref,
				},
			),
		})
		if err != nil {
			logE.WithError(err).Error("send a check to DataDog")
		}
	}

	return nil
}
