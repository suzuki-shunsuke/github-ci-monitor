AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: >
  github-ci-monitor

  Check the CI results of GitHub Repositories and send results to DataDog as Service Check.

# More info about Globals: https://github.com/awslabs/serverless-application-model/blob/master/docs/globals.rst
Globals:
  Function:
    Timeout: 5

Resources:
  GitHubCIMonitorFunction:
    Type: AWS::Serverless::Function # More info about Function Resource: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
    Properties:
      CodeUri: cmd/github-ci-monitor-lambda
      Handler: main
      Runtime: go1.x
      Tracing: Active # https://docs.aws.amazon.com/lambda/latest/dg/lambda-x-ray.html
      Environment: # More info about Env Vars: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#environment-object
        Variables:
          CONFIG: |
            repos:
            - owner: suzuki-shunsuke
              repo: test-github-action
              ref: master
              status: true
              check_run: true
              check_suite: true
            secret_id: github-ci-monitor
            region: ap-northeast-1
