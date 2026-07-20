package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string        `mapstructure:"server_address"`
	DBDriver      string        `mapstructure:"db_driver"`
	DBSource      string        `mapstructure:"db_source"`
	SymmetricKey  string        `mapstructure:"symmetrick_key"`
	Duration      time.Duration `mapstructure:"duration"`
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
