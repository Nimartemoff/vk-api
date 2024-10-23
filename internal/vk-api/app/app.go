package app

import (
	"encoding/json"
	"fmt"
	"github.com/Nimartemoff/vk-api/cmd/vk-api/config"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase/rest"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	userID   = 1
	fileName = "../durov.json"
)

func Run(cfg *config.Config) {
	c := rest.NewVKClient(cfg.VKAPI.URLs, cfg.VKAPI.Token)
	userUsecase := usecase.NewUserUsecase(c)

	user, err := userUsecase.GetUser(userID)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(user); err != nil {
		log.Error().Err(err).Send()
		return
	}
}
