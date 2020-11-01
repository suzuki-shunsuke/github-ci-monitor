package cli

import (
	"fmt"

	"github.com/suzuki-shunsuke/github-ci-monitor/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner Runner) setCLIArg(c *cli.Context, params controller.Params) controller.Params {
	if token := c.String("github-token"); token != "" {
		params.GitHubToken = token
	}
	if token := c.String("datadog-api-key"); token != "" {
		params.DataDogAPIKey = token
	}
	if logLevel := c.String("log-level"); logLevel != "" {
		params.LogLevel = logLevel
	}
	// params.ConfigFilePath = c.String("config")
	return params
}

func (runner Runner) action(c *cli.Context) error {
	params := controller.Params{}
	params = runner.setCLIArg(c, params)

	ctrl, params, err := controller.New(c.Context, params)
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}

	return ctrl.Run(c.Context, params)
}
