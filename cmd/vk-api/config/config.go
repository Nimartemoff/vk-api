package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type vkAPI struct {
	URLs  []string `env-default:"https://api.vk.com,https://api.vk.ru"`
	Token string   `env:"TOKEN" env-default:"<TOKEN HERE>"`
}

type API struct {
	Port string `env:"PORT" env-default:":8080"`
}

type neo4j struct {
	URL    string `env:"URL" env-default:"bolt://localhost:7687"`
	DBName string `env:"DB_NAME" env-default:"nizamov_vk"`
}

type Config struct {
	API            API
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
