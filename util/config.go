package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AccountName string `mapstructure:"ACCOUNT_NAME"`
	AccountKey string `mapstructure:"ACCOUNT_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	fmt.Println("Path: " + path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}