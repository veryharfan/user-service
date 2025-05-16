package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Port string    `mapstructure:"port" validate:"required"`
	Host string    `mapstructure:"host"`
	Db   DbConfig  `mapstructure:",squash"`
	Jwt  JwtConfig `mapstructure:",squash"`
}

type DbConfig struct {
	Host     string `mapstructure:"db_host" validate:"required"`
	Port     string `mapstructure:"db_port" validate:"required"`
	Username string `mapstructure:"db_username" validate:"required"`
	Password string `mapstructure:"db_password" validate:"required"`
	DbName   string `mapstructure:"db_dbname" validate:"required"`
	SSLMode  string `mapstructure:"db_sslmode"`
}

type JwtConfig struct {
	SecretKey string `mapstructure:"jwt_secretkey" validate:"required"`
	Expire    int64  `mapstructure:"jwt_expire" validate:"required"`
}

func InitConfig(ctx context.Context) (*Config, error) {
	var cfg Config

	viper.SetConfigType("env")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if !os.IsNotExist(err) {
		viper.SetConfigFile(envFile)

		if err := viper.ReadInConfig(); err != nil {
			slog.ErrorContext(ctx, "[InitConfig] ReadInConfig", "error", err)
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		slog.ErrorContext(ctx, "[InitConfig] Unmarshal", "failed bind config", err)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			slog.ErrorContext(ctx, "[InitConfig] Validation", "error", err)
		}
		slog.ErrorContext(ctx, "[InitConfig] Validation", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "[InitConfig] Config loaded", "config", cfg)
	return &cfg, nil
}
