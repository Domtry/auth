package config

import (
	"os"

	"github.com/spf13/viper"
)

type DbConfig struct {
	DbProvider string `mapstructure:"DB_PROVIDER"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbName     string `mapstructure:"DB_NAME"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
}

func LoadDBonfig() (config DbConfig, err error) {
	configFileName := "app.yaml"

	if _, err := os.Stat(configFileName); !os.IsNotExist(err) {
		viper.SetConfigFile(configFileName)
	}

	//auto loading env variable
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
