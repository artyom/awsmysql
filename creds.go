// Package awsmysql provides shortcut functions to retrieve MySQL credentials
// from AWS Secrets Manager profile.
//
// Usage:
//
//  cfg, err := awsmysql.Config(ctx, "production/dbhost")
//  if err != nil {
//      return err
//  }
//  // adjust config setting timeouts, database name, etc:
//  // see https://pkg.go.dev/github.com/go-sql-driver/mysql#Config
//  cfg.DBName = "data"
//  connector, err := mysql.NewConnector(cfg) // package github.com/go-sql-driver/mysql
//  if err != nil {
//      return err
//  }
//  db := sql.OpenDB(connector)
package awsmysql

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-sql-driver/mysql"
)

// Config is a shortcut function that creates AWS SDK session and fetches MySQL
// credentials from specific AWS Secrets Manager profile.
func Config(ctx context.Context, profile string) (*mysql.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)
	return ConfigFromSecrets(ctx, svc, profile)
}

// ConfigFromSecrets fetches MySQL credentials from specific AWS Secrets
// Manager profile using provided SecretsManager AWS SDK client.
func ConfigFromSecrets(ctx context.Context, svc *secretsmanager.Client, profile string) (*mysql.Config, error) {
	res, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &profile,
	})
	if err != nil {
		return nil, err
	}
	creds := struct {
		User string `json:"username"`
		Pass string `json:"password"`
		Host string `json:"host"`
		Port int    `json:"port"`
	}{}
	if err := json.Unmarshal([]byte(*res.SecretString), &creds); err != nil {
		return nil, err
	}
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(creds.Host, strconv.Itoa(creds.Port))
	cfg.User, cfg.Passwd = creds.User, creds.Pass
	return cfg, nil
}
