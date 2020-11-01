package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	lmb "github.com/suzuki-shunsuke/github-ci-monitor/pkg/lambda"
)

func main() {
	handler := lmb.Handler{}
	if local := os.Getenv("RUN_LOCAL"); local != "" {
		_ = handler.Start()
		return
	}
	lambda.Start(handler.Start)
}
