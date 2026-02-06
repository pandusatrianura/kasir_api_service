package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		log.Println("Local .env file found, loading config from file...")
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		_ = viper.ReadInConfig()
	} else {
		log.Println("No .env file found, using OS Environment Variables")
	}
}
