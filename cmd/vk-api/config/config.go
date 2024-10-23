package config

import "github.com/ilyakaznacheev/cleanenv"

type VKAPI struct {
	URLs  []string `env-default:"https://api.vk.com,https://api.vk.ru"`
	Token string   `env:"TOKEN" env-required:"true"`
}

type Config struct {
	VKAPI VKAPI
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}
