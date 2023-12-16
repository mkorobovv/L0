package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Nats     Nats     `yaml:"nats"`
	Postgres Postgres `yaml:"postgres"`
	App      App      `yaml:"app"`
	HTTP     HTTP     `yaml:"http"`
}
type HTTP struct {
	Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

type App struct {
	Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
	Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
}

type Postgres struct {
	Host     string `env-required:"true" yaml:"host" env:"PG_HOST"`
	Port     string `env-required:"true" yaml:"port" env:"PG_PORT"`
	User     string `env-required:"true" yaml:"user" env:"PG_USER"`
	Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
	DBName   string `env-required:"true" yaml:"name" env:"PG_NAME"`
	PgDriver string `env-required:"true" yaml:"pg_driver" env:"PG_PG_DRIVER"`
}

type Nats struct {
	Host    string `env-required:"true" yaml:"host" env:"NATS_HOST"`
	Port    string `env-required:"true" yaml:"port" env:"NATS_PORT"`
	Cluster string `env-required:"true" yaml:"cluster" env:"NATS_CLUSTER"`
	Client  string `env-required:"true" yaml:"client" env:"NATS_CLIENT"`
	Topic   string `env-required:"true" yaml:"topic" env:"NATS_TOPIC"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = &Config{}

	if err = cleanenv.ReadConfig("/Users/maxim/Documents/myprojects/L0/config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("read config error: %v", err)
	}
	if err = cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read env error: %v", err)
	}
	return cfg, nil
}
