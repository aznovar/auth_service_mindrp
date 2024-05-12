package config

import (
	"flag"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string        `mapstructure:"env"`
	GRPC     GRPCConfig    `mapstructure:"grpc"`
	DB       DBConfig      `mapstructure:"db"`
	TokenTTL time.Duration `mapstructure:"token_ttl"`
}

type GRPCConfig struct {
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// MustLoad initializes and validates config from the config file
func MustLoad() *Config {
	configPath := fetchConfigPath() // Получаем путь к конфигурационному файлу
	viper.SetConfigFile(configPath) // Указываем полный путь к файлу конфигурации

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %s", err)
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
		if res == "" {
			res = "./config/config_local.yaml"
		}
	}

	return res
}
