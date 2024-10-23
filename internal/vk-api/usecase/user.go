package usecase

import (
	"context"
	"fmt"
	"github.com/Nimartemoff/vk-api/internal/vk-api/models"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase/rest"
	"time"
)

const requestTimeout = 60 * time.Second

type UserUsecase struct {
	client *rest.VKClient
}

func NewUserUsecase(client *rest.VKClient) *UserUsecase {
	return &UserUsecase{client}
}

func (uc *UserUsecase) GetUser(userID uint64) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	user, err := uc.client.GetUsers(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("uc.client.GetUsers: %w", err)
	}

	if len(user) == 0 {
		return models.User{}, fmt.Errorf("uc.client.GetUsers: user not found")
	}

	followers, err := uc.client.GetFollowers(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("uc.client.GetFollowers: %w", err)
	}

	user[0].Followers = followers

	subscriptions, err := uc.client.GetSubscriptions(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("uc.client.GetSubscriptions: %w", err)
	}

	user[0].Subscriptions = subscriptions

	return user[0], nil
}
