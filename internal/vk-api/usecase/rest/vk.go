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
	apiVersion                 = "5.131"

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
				SetQueryParam("v", apiVersion)

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
		Count       uint64   `json:"count"`
		FollowerIDs []uint64 `json:"items"`
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

		followers = make([]models.User, 0, response.Params.Count)
		for _, followerID := range response.Params.FollowerIDs {
			followers = append(followers, models.User{ID: followerID})
		}

		return followers, nil
	}

	return nil, err
}

func (c *VKClient) GetSubscriptions(ctx context.Context, id uint64) (subscription models.Subscription, err error) {
	userID := strconv.FormatUint(id, 10)
	var resp *resty.Response

	type content struct {
		Count uint64   `json:"count"`
		Items []uint64 `json:"items"`
	}

	type respParams struct {
		Users  content `json:"users"`
		Groups content `json:"groups"`
	}

	var response struct {
		Params respParams `json:"response"`
	}

	for _, baseURL := range c.baseURLs {
		urlr, err := url.JoinPath(baseURL, getSubscriptionsMethodName)
		if err != nil {
			return subscription, err
		}

		if err := rest.GetRestClient(func() (errGet error) {
			req := c.resty.R().
				SetContext(ctx).
				SetQueryParams(map[string]string{
					"v":       apiVersion,
					"user_id": userID,
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

		subscription.Users = make([]models.User, 0, response.Params.Users.Count)
		for _, user := range response.Params.Users.Items {
			subscription.Users = append(subscription.Users, models.User{ID: user})
		}

		subscription.Groups = make([]models.Group, 0, response.Params.Groups.Count)
		for _, user := range response.Params.Groups.Items {
			subscription.Groups = append(subscription.Groups, models.Group{ID: user})
		}

		return subscription, nil
	}

	return subscription, err
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
