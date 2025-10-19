package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	PrimaryEnv           string   `koanf:"primary_env" validate:"required"`
	ServerPort           string   `koanf:"server_port" validate:"required"`
	ServerReadTimeout    int      `koanf:"server_read_timeout" validate:"required"`
	ServerWriteTimeout   int      `koanf:"server_write_timeout" validate:"required"`
	ServerIdleTimeout    int      `koanf:"server_idle_timeout" validate:"required"`
	ServerCORSAllowedOrigins []string `koanf:"server_cors_allowed_origins" validate:"required"`

	DatabaseHost         string `koanf:"database_host" validate:"required"`
	DatabasePort         int    `koanf:"database_port" validate:"required"`
	DatabaseUser         string `koanf:"database_user" validate:"required"`
	DatabasePassword     string `koanf:"database_password"`
	DatabaseName         string `koanf:"database_name" validate:"required"`
	DatabaseSSLMode      string `koanf:"database_ssl_mode" validate:"required"`
	DatabaseMaxOpenConns  int    `koanf:"database_max_open_conns" validate:"required"`
	DatabaseMaxIdleConns  int    `koanf:"database_max_idle_conns" validate:"required"`
	DatabaseConnMaxLifetime int `koanf:"database_conn_max_lifetime" validate:"required"`
	DatabaseConnMaxIdleTime int `koanf:"database_conn_max_idle_time" validate:"required"`

	RedisAddress         string `koanf:"redis_address" validate:"required"`

	AuthSecretKey        string `koanf:"auth_secret_key" validate:"required"`

	IntegrationResendAPIKey string `koanf:"integration_resend_api_key" validate:"required"`

	Observability *ObservabilityConfig `koanf:"observability"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	err := k.Load(env.Provider("BOILERPLATE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "BOILERPLATE_"))
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	// Set default observability config if not provided
	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	// Override service name and environment from primary config
	mainConfig.Observability.ServiceName = "boilerplate"
	mainConfig.Observability.Environment = mainConfig.PrimaryEnv

	// Validate observability config
	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	return mainConfig, nil
}
