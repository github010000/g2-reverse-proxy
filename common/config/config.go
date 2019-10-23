package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type config struct {
	Env         string
	ListenPort  int
	OpsgenieURL string
}

var _config *config

func LoadConfig() error {
	_config = &config{}

	if err := loadEnvAsStr(&_config.Env, "ENV", true); err != nil {
		return err
	}
	if err := loadEnvAsInt(&_config.ListenPort, "LISTEN_PORT", true); err != nil {
		return err
	}
	if err := loadEnvAsStr(&_config.OpsgenieURL, "OPSGENIE_URL", true); err != nil {
		return err
	}
	return nil
}

func loadEnvAsInt(configVal *int, envKey string, isRequired bool) error {
	envVal, err := strconv.Atoi(os.Getenv(envKey))
	if err != nil && isRequired == true {
		return fmt.Errorf("ENV %s required", envKey)
	}

	if err == nil {
		*configVal = envVal
	}
	return nil
}

func loadEnvAsStr(configVal *string, envKey string, isRequired bool) error {
	envVal := os.Getenv(envKey)
	if envVal == "" && isRequired == true {
		return fmt.Errorf("ENV %s required", envKey)
	}

	if envVal != "" {
		*configVal = strings.TrimSpace(envVal)
	}

	return nil
}

func Env() string {
	return _config.Env
}

func IsDev() bool {
	return Env() == "dev"
}

func IsProd() bool {
	return !IsDev()
}

func OpsgenieURL() string {
	return _config.OpsgenieURL
}
