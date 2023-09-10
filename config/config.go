// package config

// import (
// 	"fmt"

// 	"github.com/ilyakaznacheev/cleanenv"
// )

// type (
// 	Config struct {
// 		App      `yaml:"app"`
// 		HTTP     `yaml:"http"`
// 		Postgres `yaml:"postgres`
// 	}

// 	App struct {
// 		Name    string `yaml:"name" env:"APPLICATION_NAME"`
// 		Version string `yaml:"version" env:"APPLICATION_VERSION"`
// 	}

// 	HTTP struct {
// 		Host string `yaml:"host" env:"HOST"`
// 		Port string `yaml:"port" env:"PORT"`
// 	}

// 	Postgres struct {
// 		User         string `yaml:"user" env:"USER"`
// 		Password     string `yaml:"password" env:"PASSWORD"`
// 		DatabaseName string `yaml:"database_name" env:"DATABASE_NAME"`
// 	}
// )

// func NewConfig() (*Config, error) {
// 	cfg := &Config{}

// 	err := cleanenv.ReadConfig("./config/config.yml", cfg)
// 	if err != nil {
// 		return nil, fmt.Errorf("config error: %v", err)
// 	}

// 	err = cleanenv.ReadEnv(cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return cfg, nil
// }
