package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HttpServerAddress    string        `mapstructure:"http_server_address"`
	GrpcServerAddress    string        `mapstructure:"grpc_server_address"`
	DBDriver             string        `mapstructure:"db_driver"`
	DBSource             string        `mapstructure:"db_source"`
	SymmetricKey         string        `mapstructure:"symmetrick_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return
}
