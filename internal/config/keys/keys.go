package keys

import (
	"github.com/spf13/viper"
)

// Some commonly used configuration keys for Viper.
var (
	ConfigKeyEnv = key[string]("eelchat.deploy.env")

	ConfigKeyTlsKey    = key[string]("eelchat.server.tlsKey")
	ConfigKeyTlsCert   = key[string]("eelchat.server.tlsCert")
	ConfigKeyServeHost = key[string]("eelchat.server.host")
	ConfigKeyServePort = key[string]("eelchat.server.port")

	ConfigKeyDbName = key[string]("eelchat.datasource.dbname")
	ConfigKeyDbUser = key[string]("eelchat.datasource.user")
	ConfigKeyDbPass = key[string]("eelchat.datasource.pass")
	ConfigKeyDbHost = key[string]("eelchat.datasource.host")
	ConfigKeyDbPort = key[string]("eelchat.datasource.port")
)

type ConfigKey[T any] struct {
	Key string
}

func (c ConfigKey[T]) Get() T {
	return viper.Get(c.Key).(T)
}

func key[T any](key string) ConfigKey[T] {
	return ConfigKey[T]{Key: key}
}
