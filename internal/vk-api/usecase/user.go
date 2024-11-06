package usecase

import (
	"context"
	"fmt"
	"github.com/Nimartemoff/vk-api/internal/vk-api/models"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase/repo/neo4j"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase/rest"
	"github.com/rs/zerolog/log"
	"time"
)

const requestTimeout = 60 * time.Second

type UserUsecase struct {
	client    *rest.VKClient
	neo4jRepo *neo4j.UserNeo4jRepo
}

func NewUserUsecase(client *rest.VKClient, neo4jRepo *neo4j.UserNeo4jRepo) *UserUsecase {
	return &UserUsecase{client: client, neo4jRepo: neo4jRepo}
}

func (uc *UserUsecase) CreateIndexes(ctx context.Context) error {
	return uc.neo4jRepo.CreateIndexes(ctx)
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

	log.Info().Msgf("Обошел пользователя: %s %s", user[0].FirstName, user[0].LastName)
	return user[0], nil
}

func (uc *UserUsecase) GetUsersWithDepth(userID uint64, depth int) (models.User, error) {
	if depth <= 0 {
		return models.User{}, nil
	}

	user, err := uc.GetUser(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("uc.GetUser: %w", err)
	}

	for i := range user.Followers {
		if user.Followers[i], err = uc.GetUsersWithDepth(user.Followers[i].ID, depth-1); err != nil {
			return models.User{}, fmt.Errorf("uc.GetUsersWithDepth: %w", err)
		}
	}

	for i := range user.Subscriptions.Users {
		if user.Subscriptions.Users[i], err = uc.GetUsersWithDepth(user.Subscriptions.Users[i].ID, depth-1); err != nil {
			return models.User{}, fmt.Errorf("uc.GetUsersWithDepth: %w", err)
		}
	}

	return user, nil
}

func (uc *UserUsecase) SaveUser(ctx context.Context, user models.User) error {
	if user.ID == 0 {
		return nil
	}

	log.Info().Msgf("Создание пользователя %s %s", user.FirstName, user.LastName)
	if err := uc.neo4jRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	for _, follower := range user.Followers {
		if follower.ID == 0 {
			continue
		}

		if err := uc.SaveUser(ctx, follower); err != nil {
			return err
		}

		log.Info().Msgf("Создание отношения (%s %s)->[:FOLLOW]->(%s %s)", follower.FirstName, follower.LastName, user.FirstName, user.LastName)
		if err := uc.neo4jRepo.CreateFollowRelationship(ctx, follower, user); err != nil {
			return err
		}
	}

	for _, subscription := range user.Subscriptions.Users {
		if subscription.ID == 0 {
			continue
		}

		if err := uc.SaveUser(ctx, subscription); err != nil {
			return err
		}

		log.Info().Msgf("Создание отношения (%s %s)->[:SUBSCRIBE]->(%s %s)", user.FirstName, user.LastName, subscription.FirstName, subscription.LastName)
		if err := uc.neo4jRepo.CreateSubscribeUserUserRelationship(ctx, user, subscription); err != nil {
			return err
		}
	}

	for _, group := range user.Subscriptions.Groups {
		if err := uc.neo4jRepo.CreateGroup(ctx, group); err != nil {
			return err
		}

		log.Info().Msgf("Создание отношения (%s %s)->[:SUBSCRIBE]->(%s)", user.FirstName, user.LastName, group.Name)
		if err := uc.neo4jRepo.CreateSubscribeUserGroupRelationship(ctx, user, group); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UserUsecase) GetUsersCount(ctx context.Context) (int, error) {
	return uc.neo4jRepo.GetUsersCount(ctx)
}

func (uc *UserUsecase) GetGroupsCount(ctx context.Context) (int, error) {
	return uc.neo4jRepo.GetGroupsCount(ctx)
}

func (uc *UserUsecase) GetTopUsersByFollowersCount(ctx context.Context, limit int) ([]models.User, error) {
	return uc.neo4jRepo.GetTopUsersByFollowersCount(ctx, limit)
}

func (uc *UserUsecase) GetTopGroupsBySubscribersCount(ctx context.Context, limit int) ([]models.Group, error) {
	return uc.neo4jRepo.GetTopGroupsBySubscribersCount(ctx, limit)
}

func (uc *UserUsecase) GetUsersWithDifferentGroups(ctx context.Context, limit int) ([]models.User, error) {
	return uc.neo4jRepo.GetUsersWithDifferentGroups(ctx, limit)
}
