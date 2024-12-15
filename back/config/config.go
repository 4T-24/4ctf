package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var (
	EnvironmentDevelopment = "dev"
	EnvironmentProduction  = "prod"
)

type Config struct {
	Server struct {
		Environment string `mapstructure:"environment"`
		Port        int    `mapstructure:"port"`

		Key string `mapstructure:"key"`
	} `mapstructure:"server"`

	MySql struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`

		MigrationsFolder string `mapstructure:"migrations_folder"`
		SchemaFile       string `mapstructure:"schema_file"`
	} `mapstructure:"mysql"`
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == EnvironmentDevelopment
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == EnvironmentProduction
}

func listStructKeys(s interface{}) ([]string, error) {
	// Recursively get the config struct tag mapstructure
	keys := []string{}
	ct := reflect.TypeOf(s)

	if ct.Kind() != reflect.Struct {
		return nil, fmt.Errorf("listStructKeys: %v is not a struct", ct)
	}

	for i := range ct.NumField() {
		field := ct.Field(i)
		tag := field.Tag.Get("mapstructure")

		if field.Type.Kind() == reflect.Struct {
			res, err := listStructKeys(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return nil, err
			}
			for _, k := range res {
				keys = append(keys, fmt.Sprintf("%s.%s", tag, k))
			}
		} else {
			keys = append(keys, tag)
		}
	}

	return keys, nil
}

func bindAllEnv(v *viper.Viper, s interface{}) error {
	keys, err := listStructKeys(s)
	if err != nil {
		return err
	}

	for _, k := range keys {
		v.BindEnv(k)
	}

	return nil
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

	// Read config
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, using defaults: %v", err)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindAllEnv(v, Config{})

	// Unmarshal into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	return &config
}
