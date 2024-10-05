package config

import (
	"github.com/spf13/viper"
	"log"
)

type ApiConfig struct {
	ApiAddress string `mapstructure:"API_ADDRESS"`
	DbDrives   string `mapstructure:"DB_DRIVER"`
	DbAddress  string `mapstructure:"DB_ADDRESS"`
}

func NewApiConfig() *ApiConfig {
	return &ApiConfig{}
}

func (config *ApiConfig) Parse() {
	viper.SetConfigName(".env")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("Config not found")
		}
		log.Fatalln(err)
	}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
