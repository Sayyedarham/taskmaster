package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secrets struct {
	DatabaseURL   string `json:"database_url"`
	RedisPassword string `json:"redis_password"`
	JWTSecret     string `json:"jwt_secret"`
}

func LoadSecretsFromAWS(secretName string) (*Secrets, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return nil, err
	}

	var secrets Secrets
	if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		return nil, err
	}

	return &secrets, nil
}

func LoadConfigWithSecrets() *Config {
	base := Load()

	// In production, override with AWS Secrets Manager
	if os.Getenv("ENV") == "production" {
		secretName := os.Getenv("AWS_SECRET_NAME")
		if secretName != "" {
			secrets, err := LoadSecretsFromAWS(secretName)
			if err == nil {
				base.DatabaseURL = secrets.DatabaseURL
				base.RedisPassword = secrets.RedisPassword
				base.JWTSecret = secrets.JWTSecret
			}
		}
	}

	return base
}
