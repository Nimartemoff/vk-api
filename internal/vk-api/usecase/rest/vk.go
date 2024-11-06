package rest

import (
	"context"
	"github.com/Nimartemoff/vk-api/internal/vk-api/models"
	"github.com/Nimartemoff/vk-api/pkg/rest"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	clientRetryCount       = 5
	clientRetryWaitTime    = time.Second * 1
	clientRetryMaxWaitTime = time.Second * 5
	correctionTime         = 15 * time.Second

	usersGetMethodName         = "method/users.get"
	getFollowersMethodName     = "method/users.getFollowers"
	getSubscriptionsMethodName = "method/users.getSubscriptions"
	apiVersion                 = "5.199"

	usersGetFields         = "screen_name,sex,city"
	getSubscriptionsFields = "name,screen_name"

	count = 3

	base = 10
)

type VKClient struct {
	baseURLs []string
	resty    *resty.Client
}

func NewVKClient(baseURLs []string, token string) *VKClient {
	rc := resty.New()
	rc.SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime).
		AddRetryAfterErrorCondition().
		SetAuthToken(token).
		SetQueryParam("lang", "ru")

	return &VKClient{
		baseURLs: baseURLs,
		resty:    rc,
	}
}

func (c *VKClient) GetUsers(ctx context.Context, ids ...uint64) ([]models.User, error) {
	userIDs, err := c.parseUserIDs(ids...)
	if err != nil {
		return nil, err
	}

	var resp *resty.Response

	var response struct {
		Users []models.User `json:"response"`
	}

	for _, baseURL := range c.baseURLs {
		urlr, err := url.JoinPath(baseURL, usersGetMethodName)
		if err != nil {
			return nil, err
		}

		if err := rest.GetRestClient(func() (errGet error) {
			req := c.resty.R().
				SetContext(ctx).
				SetQueryParams(map[string]string{
					"v":      apiVersion,
					"fields": usersGetFields,
				})

			if userIDs != "" {
				req.SetQueryParam("user_ids", userIDs)
			}

			resp, errGet = req.Get(urlr)

			if errGet != nil {
				return errGet
			}

			return nil
		}); err != nil {
			log.Error().Err(err).Send()
			continue
		}

		body := resp.Body()
		if body == nil {
			log.Warn().Msg("body == nil")
			continue
		}

		if err := sonic.Unmarshal(body, &response); err != nil {
			log.Error().Err(err).Send()
			continue
		}

		return response.Users, nil
	}

	return nil, err
}

func (c *VKClient) GetFollowers(ctx context.Context, id uint64) (followers []models.User, err error) {
	userID := strconv.FormatUint(id, 10)
	var resp *resty.Response

	type respParams struct {
		Count     uint64        `json:"count"`
		Followers []models.User `json:"items"`
	}
	var response struct {
		Params respParams `json:"response"`
	}

	for _, baseURL := range c.baseURLs {
		urlr, err := url.JoinPath(baseURL, getFollowersMethodName)
		if err != nil {
			return nil, err
		}

		if err := rest.GetRestClient(func() (errGet error) {
			req := c.resty.R().
				SetContext(ctx).
				SetQueryParams(map[string]string{
					"v":       apiVersion,
					"user_id": userID,
					"count":   strconv.FormatUint(count, 10),
					"fields":  usersGetFields,
				})

			resp, errGet = req.Get(urlr)

			if errGet != nil {
				return errGet
			}

			return nil
		}); err != nil {
			log.Error().Err(err).Send()
			continue
		}

		body := resp.Body()
		if body == nil {
			log.Warn().Msg("body == nil")
			continue
		}

		if err := sonic.Unmarshal(body, &response); err != nil {
			log.Error().Err(err).Send()
			continue
		}

		return response.Params.Followers, nil
	}

	return nil, err
}

func (c *VKClient) GetSubscriptions(ctx context.Context, id uint64) (subscriptions models.Subscriptions, err error) {
	userID := strconv.FormatUint(id, 10)
	var resp *resty.Response

	type userGroup struct {
		ID         uint64      `json:"id"`
		Type       string      `json:"type"`
		Name       string      `json:"name"`
		ScreenName string      `json:"screen_name"`
		FirstName  string      `json:"first_name"`
		LastName   string      `json:"last_name"`
		Sex        byte        `json:"sex"`
		City       models.City `json:"city"`
	}
	type respParams struct {
		Count         uint64      `json:"count"`
		Subscriptions []userGroup `json:"items"`
	}
	var response struct {
		Params respParams `json:"response"`
	}

	for _, baseURL := range c.baseURLs {
		urlr, err := url.JoinPath(baseURL, getSubscriptionsMethodName)
		if err != nil {
			return subscriptions, err
		}

		if err := rest.GetRestClient(func() (errGet error) {
			req := c.resty.R().
				SetContext(ctx).
				SetQueryParams(map[string]string{
					"v":        apiVersion,
					"user_id":  userID,
					"extended": "1",
					"count":    strconv.FormatUint(count, 10),
				})

			resp, errGet = req.Get(urlr)

			if errGet != nil {
				return errGet
			}

			return nil
		}); err != nil {
			log.Error().Err(err).Send()
			continue
		}

		body := resp.Body()
		if body == nil {
			log.Warn().Msg("c.GetUsers: body == nil")
			continue
		}

		if err := sonic.Unmarshal(body, &response); err != nil {
			log.Error().Err(err).Send()
			continue
		}

		for _, userGroup := range response.Params.Subscriptions {
			switch userGroup.Type {
			case "profile":
				subscriptions.Users = append(subscriptions.Users, models.User{
					ID:         userGroup.ID,
					ScreenName: userGroup.ScreenName,
					FirstName:  userGroup.FirstName,
					LastName:   userGroup.LastName,
					Sex:        userGroup.Sex,
					City:       userGroup.City,
				})
			case "page":
				subscriptions.Groups = append(subscriptions.Groups, models.Group{
					ID:         userGroup.ID,
					Name:       userGroup.Name,
					ScreenName: userGroup.ScreenName,
				})
			}
		}

		log.Info().Msgf("Обошел подписки: %+v", subscriptions)
		return subscriptions, nil
	}

	return subscriptions, err
}

func (c *VKClient) parseUserIDs(ids ...uint64) (string, error) {
	var builder strings.Builder

	for i, id := range ids {
		if i > 0 {
			builder.WriteString(",")
		}

		builder.WriteString(strconv.FormatUint(id, base))
	}

	return builder.String(), nil
}
