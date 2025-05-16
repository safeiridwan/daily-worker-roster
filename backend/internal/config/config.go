package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultServerPort         = 5000
	defaultJWTExpirationHours = 72
)

type Config struct {
	ServerPort    int    `yaml:"server_port" env:"SERVER_PORT"`
	DSN           string `yaml:"dsn" env:"DSN"`
	JWTExpiration int    `yaml:"jwt_expiration" env:"JWT_EXPIRATION"`
	JWTSigningKey string `yaml:"jwt_signing_key" env:"JWT_SIGNING_KEY"`
}

func (c Config) Validate() error {
	validate := validator.New()
	rules := map[string]string{
		"DSN": "required",
	}
	validate.RegisterStructValidationMapRules(rules, Config{})
	return validate.Struct(c)
}

func Load(file string, logger zerolog.Logger) (*Config, error) {
	c := Config{
		ServerPort: defaultServerPort,
	}

	bytes, err := os.ReadFile(file)

	if err != nil {
		logger.Error().Err(err)
		return nil, err
	}

	if err = yaml.Unmarshal(bytes, &c); err != nil {
		logger.Error().Err(err)
		return nil, err
	}

	if err = c.Validate(); err != nil {
		logger.Error().Err(err)
		return nil, err
	}

	return &c, err
}
