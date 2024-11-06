package main

import (
	"github.com/Nimartemoff/vk-api/cmd/vk-api/config"
	"github.com/Nimartemoff/vk-api/internal/vk-api/app"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read env")
	}

	// Run
	app.Run(cfg)
}
