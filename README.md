# github-ci-monitor

[![Build Status](https://github.com/suzuki-shunsuke/github-ci-monitor/workflows/CI/badge.svg)](https://github.com/suzuki-shunsuke/github-ci-monitor/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/github-ci-monitor)](https://goreportcard.com/report/github.com/suzuki-shunsuke/github-ci-monitor)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/github-ci-monitor.svg)](https://github.com/suzuki-shunsuke/github-ci-monitor)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/github-ci-monitor/master/LICENSE)

Check GitHub repositories CI statues and send them to DataDog.

## Overview

When we merge pull requests, sometimes CI/CD fails.
Then we want to find the failure.

By running this tool periodically, we can monitor the status by DataDog.

This tool gets GitHub repositories CI statues (commit status and Check API) and sends the result by DataDog API.

## How to run this tool as Lambda Function

* Create AWS Secrets Manager's secret
* Deploy the Lambda Function with AWS SAM
* Give the IAM Role the permission to read the secret
* Configure AWS CloudWatch Events to run this tool periodically

### Configuration

* Lambda's environment variables
* AWS Secrets Manager's secret

#### Lambda's environment variables

* `CONFIG`: YAML

```yaml
---
repos:
- owner: suzuki-shunsuke
  repo: test-github-action
  ref: master
  status: true
  check_run: true
  check_suite: true
  tags:
    codeowner: sre
region: ap-northeast-1
secret_id: github-ci-monitor
check_name: ci.ok
tags:
  sender: github-ci-monitor
# version_id:
```

#### AWS Secrets Manager's secret

* datadog_api_key
* github_token

## How to create the DataDog Monitor

New Monitor > Custom Check

* Select a check: `ci.ok`

## LICENSE

[MIT](LICENSE)
