package lambda

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/github-ci-monitor/pkg/controller"
	"gopkg.in/yaml.v3"
)

type Handler struct{}

type Config struct {
	Repos     []controller.Repo
	LogLevel  string `yaml:"log_level"`
	SecretID  string `yaml:"secret_id"`
	VersionID string `yaml:"version_id"`
	Region    string
}

type Secret struct {
	GitHubToken   string `yaml:"github_token"`
	DataDogAPIKey string `yaml:"datadog_api_key"`
}

func (handler Handler) start() error {
	cfgString := os.Getenv("CONFIG")
	cfg := Config{}
	if cfgString != "" {
		if err := yaml.Unmarshal([]byte(cfgString), &cfg); err != nil {
			return fmt.Errorf("unmarshal config: %w", err)
		}
	}
	ctx := context.Background()

	sess := session.Must(session.NewSession())

	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(cfg.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(cfg.SecretID),
	}
	if cfg.VersionID != "" {
		input.VersionId = aws.String(cfg.VersionID)
	}
	output, err := svc.GetSecretValueWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("get secret value from AWS SecretsManager: %w", err)
	}
	secret := Secret{}
	if err := yaml.Unmarshal([]byte(*output.SecretString), &secret); err != nil {
		return fmt.Errorf("parse secret value: %w", err)
	}
	params := controller.Params{
		Repos:         cfg.Repos,
		LogLevel:      cfg.LogLevel,
		GitHubToken:   secret.GitHubToken,
		DataDogAPIKey: secret.DataDogAPIKey,
	}
	ctrl, params, err := controller.New(ctx, params)
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}

	return ctrl.Run(ctx, params)
}

func (handler Handler) Start() error {
	if err := handler.start(); err != nil {
		logrus.WithError(err).Error("start")
		return err
	}
	return nil
}
