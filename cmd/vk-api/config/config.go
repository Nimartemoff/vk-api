package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type vkAPI struct {
	URLs  []string `env-default:"https://api.vk.com,https://api.vk.ru"`
	Token string   `env:"TOKEN" env-required:"true"`
}

type neo4j struct {
	URL    string `env:"URL" env-required:"true"`
	DBName string `env:"DB_NAME" env-default:"nizamov_vk"`
}

type Config struct {
	VKAPI          vkAPI
	Neo4j          neo4j
	ContextTimeout time.Duration `env:"TIMEOUT" env-default:"60s"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}
