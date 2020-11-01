package cli

import (
	"context"
	"io"

	"github.com/suzuki-shunsuke/github-ci-monitor/pkg/constant"
	"github.com/urfave/cli/v2"
)

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (runner Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "github-ci-monitor",
		Usage:   "get GitHub repositories statuses and send them to DataDog. https://github.com/suzuki-shunsuke/github-ci-monitor",
		Version: constant.Version,
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "get CI information",
				Action: runner.action,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "github-token",
						Usage: "GitHub Access Token [$GITHUB_TOKEN, $GITHUB_ACCESS_TOKEN]",
					},
					&cli.StringFlag{
						Name:  "datadog-api-key",
						Usage: "DataDog API Key [$DATADOG_API_KEY]",
					},
					&cli.StringFlag{
						Name:  "log-level",
						Usage: "log level",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "configuration file path",
					},
				},
			},
		},
	}

	return app.RunContext(ctx, args)
}
