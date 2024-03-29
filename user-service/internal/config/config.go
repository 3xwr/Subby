package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Port         int    `env:"PORT" envDefault:"8080"`
	DBConnString string `env:"DB_CONN_STRING,notEmpty,unset,file"`
	Secret       string `env:"SECRET,notEmpty,unset,file"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
