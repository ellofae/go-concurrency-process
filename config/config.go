package config

import (
	"os"

	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresDB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		SSLmode  string `yaml:"sslmode"`
		MaxConns string `yaml:"maxconns"`
	}

	Server struct {
		BindAddr     string `yaml:"bindAddr"`
		ReadTimeout  string `yaml:"readTimeout"`
		WriteTimeout string `yaml:"writeTimeout"`
		IdleTimeout  string `yaml:"idleTimeout"`
	}
}

func ConfigureViper() *viper.Viper {
	logger := logger.GetLogger()

	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		logger.Error("Unable to read the configuration file.", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Config loaded successfully.")

	return v
}

func ParseConfig(v *viper.Viper) *Config {
	logger := logger.GetLogger()

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		logger.Error("Unable to parse the configuration file.")
	}
	logger.Info("Configuratin file parsed successfully.")

	return cfg
}
