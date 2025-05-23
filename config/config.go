package config

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RunAddress     string `env:"RUN_ADDRESS" env-default:"localhost:8081"`
	RootPath       string `env:"ROOT_PATH" env-default:"/api/user"`
	DatabaseURI    string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Debug          bool   `env:"DEBUG" env-default:"true"`
	AuthCookieName string `env:"AUTH_COOKIE_NAME" env-default:"GOPHERMART_ACCESS_TOKEN"`
}

var cfg Config

func initFromCLI() {
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "addr:port the service will run at")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "postgres DSN")
	flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "accrual system address")
	flag.StringVar(&cfg.RootPath, "p", cfg.RootPath, "endpoint prefix")

	flag.Parse()
}

func initFromEnv() error {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	return nil
}

func InitConfig() (*Config, error) {
	err := initFromEnv()
	initFromCLI()
	return &cfg, err
}
