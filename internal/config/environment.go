package config

import (
	"strings"

	. "github.com/apachejuice/eelchat/internal/config/keys"
)

type Environment int

func (e Environment) String() string {
	switch e {
	case EnvProduction:
		return "production"
	case EnvDevelopment:
		return "development"
	case EnvTesting:
		return "testing"
	default:
		return ""
	}
}

// Environment configuration. The environment is more of a general set of controls over some configuration options.
const (
	EnvProduction Environment = iota
	EnvDevelopment
	EnvTesting
)

var (
	env    Environment
	envMap map[string]Environment = map[string]Environment{
		"prod":       EnvProduction,
		"production": EnvProduction,

		"dev":         EnvDevelopment,
		"develop":     EnvDevelopment,
		"development": EnvDevelopment,

		"test":    EnvTesting,
		"testing": EnvTesting,
	}
)

func setEnv() {
	name := ConfigKeyEnv.Get()
	if name == "" {
		configLogger.Fatal("Unable to load environment configuration: key `" + ConfigKeyEnv.Key + "` not defined")
	}

	e, ok := envMap[strings.ToLower(name)]
	if ok {
		env = e
	} else {
		configLogger.Fatal("Unable to load environment configuration: invalid env value", ConfigKeyEnv, name)
	}
}

func GetEnv() Environment {
	return env
}
