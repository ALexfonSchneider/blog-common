package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
)

type AppEnv string

const (
	Local      AppEnv = "local"
	Docker     AppEnv = "docker"
	Production AppEnv = "production"
)

func MustLoadConfig[T any]() *T {
	selectedEnv := Local
	if os.Getenv("APP_ENV") != "" {
		selectedEnv = AppEnv(os.Getenv("APP_ENV"))
	}

	var cfg T
	k := koanf.New(".")
	configPath := fmt.Sprintf("config/%s.yaml", selectedEnv)

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(cfg); err != nil {
		panic(fmt.Errorf("config validation failed: %w", err))
	}

	return &cfg
}
