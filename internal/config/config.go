package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfig *envConfig

func InitEnvConfigs() {
	EnvConfig = loadEnvVariables()
}

type envConfig struct {
	LocalServerPort string `mapstructure:"LOCAL_SERVER_PORT"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUsername      string `mapstructure:"DB_USER_NAME"`
	DBName          string `mapstructure:"DB_NAME"`
	DBSSLMode       string `mapstructure:"DB_SSL_MODE"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
}

func loadEnvVariables() (config *envConfig) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
