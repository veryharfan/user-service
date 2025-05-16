package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Port string    `mapstructure:"PORT" validate:"required"`
	Host string    `mapstructure:"HOST"`
	Db   DbConfig  `mapstructure:",squash"`
	Jwt  JwtConfig `mapstructure:",squash"`
}

type DbConfig struct {
	Host     string `mapstructure:"DB_HOST" validate:"required"`
	Port     string `mapstructure:"DB_PORT" validate:"required"`
	Username string `mapstructure:"DB_USERNAME" validate:"required"`
	Password string `mapstructure:"DB_PASSWORD" validate:"required"`
	DbName   string `mapstructure:"DB_DBNAME" validate:"required"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type JwtConfig struct {
	SecretKey string `mapstructure:"JWT_SECRETKEY" validate:"required"`
	Expire    int64  `mapstructure:"JWT_EXPIRE" validate:"required"`
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

	fmt.Println("JWT_SECRETKEY (viper):", viper.GetString("JWT_SECRETKEY"))
	fmt.Println("os.Getenv:", os.Getenv("JWT_SECRETKEY"))

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
