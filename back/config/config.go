package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig() *Config {
	v := viper.New()

	// Add multiple config paths
	v.SetConfigName("config")
	v.AddConfigPath(".")        // Current working directory
	v.AddConfigPath("./config") // Specific directory in working dir
	v.AddConfigPath("./back/config")
	v.SetConfigType("yaml")

	// Add home directory as a config source
	home, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(home)
	}

	v.AutomaticEnv() // Allow environment variable overrides

	// Read config
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, using defaults: %v", err)
	}

	// Unmarshal into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	return &config
}
