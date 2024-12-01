package test

import (
	"fmt"
	"github.com/Nimartemoff/vk-api/internal/vk-api/models"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"strings"
	"testing"
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

	baseURL = "http://localhost:8080/api/v1/nodes"
	token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJlZGl0b3IiXSwiaWF0IjoxMjg2Njk0MDAwfQ.ND_sy9LK9j4K-vQPIVMQ9mjdPoIp8a8C_UtdWxlpvtw"
)

func TestGetNodes(t *testing.T) {
	_, err := getNodes()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCreateUser(t *testing.T) {
	user := models.User{
		ID:         1234,
		ScreenName: "NimartemX13",
		FirstName:  "Artem",
		LastName:   "Nizamov",
		Sex:        2,
		City: models.City{
			Title: "Nefteyugansk",
		},
		Followers: make([]models.User, 0),
		Subscriptions: models.Subscriptions{
			Users: []models.User{
				{
					ID:         12345,
					ScreenName: "Puritanin",
					FirstName:  "Pavel",
					LastName:   "Demukhametov",
					Sex:        2,
					City: models.City{
						Title: "Nizhnevartovsk",
					},
				},
				{
					ID:         123456,
					ScreenName: "PostalDude",
					FirstName:  "Yakov",
					LastName:   "Zarembo",
					Sex:        2,
					City: models.City{
						Title: "Tyumen",
					},
				},
			},
			Groups: nil,
		},
	}

	resp, err := resty.New().
		SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime).
		AddRetryAfterErrorCondition().
		SetAuthToken(token).
		SetBaseURL(baseURL).
		R().SetQueryParam("type", "user").SetBody(user).Post("")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		t.Fatalf("%v", resp.Status())
	}

	body := resp.Body()
	if body == nil {
		t.Fatalf("body == nil")
	}
}

func TestCreateGroup(t *testing.T) {
	group := models.GroupWithSubscribers{
		Group: models.Group{
			ID:         123456789,
			Name:       "Tyumen State University",
			ScreenName: "TSU",
		},
		Subscribers: []models.User{
			{
				ID:         12345,
				ScreenName: "Puritanin",
				FirstName:  "Pavel",
				LastName:   "Demukhametov",
				Sex:        2,
				City: models.City{
					Title: "Nizhnevartovsk",
				},
			},
			{
				ID:         1234,
				ScreenName: "NimartemX13",
				FirstName:  "Artem",
				LastName:   "Nizamov",
				Sex:        2,
				City: models.City{
					Title: "Nefteyugansk",
				},
			},
		},
	}

	resp, err := resty.New().
		SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime).
		AddRetryAfterErrorCondition().
		SetAuthToken(token).
		SetBaseURL(baseURL).
		R().SetQueryParam("type", "group").SetBody(group).Post("")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		t.Fatalf("%v", resp.Status())
	}

	body := resp.Body()
	if body == nil {
		t.Fatalf("body == nil")
	}
}

func TestGetUser(t *testing.T) {
	user, err := getUser(0)
	if err != nil {
		t.Fatalf("%v", err)
	}

	require.Equal(t, user, models.User{
		ID:         1234,
		ScreenName: "NimartemX13",
		FirstName:  "Artem",
		LastName:   "Nizamov",
		Sex:        2,
		City: models.City{
			Title: "Nefteyugansk",
		},
		Followers: nil,
		Subscriptions: models.Subscriptions{
			Users: []models.User{
				{
					ID:         12345,
					ScreenName: "Puritanin",
					FirstName:  "Pavel",
					LastName:   "Demukhametov",
					Sex:        2,
					City: models.City{
						Title: "Nizhnevartovsk",
					},
				},
				{
					ID:         123456,
					ScreenName: "PostalDude",
					FirstName:  "Yakov",
					LastName:   "Zarembo",
					Sex:        2,
					City: models.City{
						Title: "Tyumen",
					},
				},
			},
			Groups: []models.Group{{
				ID:         123456789,
				Name:       "Tyumen State University",
				ScreenName: "TSU",
			}},
		},
	})
}

func TestGetGroup(t *testing.T) {
	group, err := getGroup(3)
	if err != nil {
		t.Fatalf("%v", err)
	}

	require.Equal(t, group, models.GroupWithSubscribers{
		Group: models.Group{
			ID:         123456789,
			Name:       "Tyumen State University",
			ScreenName: "TSU",
		},
		Subscribers: []models.User{
			{
				ID:         12345,
				ScreenName: "Puritanin",
				FirstName:  "Pavel",
				LastName:   "Demukhametov",
				Sex:        2,
				City: models.City{
					Title: "Nizhnevartovsk",
				},
			},
			{
				ID:         1234,
				ScreenName: "NimartemX13",
				FirstName:  "Artem",
				LastName:   "Nizamov",
				Sex:        2,
				City: models.City{
					Title: "Nefteyugansk",
				},
			},
		},
	})
}

func TestDeleteGroup(t *testing.T) {
	if err := deleteNode(123456789); err != nil {
		t.Fatalf("%v", err)
	}

	_, err := getGroup(3)
	if err == nil || !strings.Contains(err.Error(), http.StatusText(http.StatusNotFound)) {
		t.Fatalf("%v", err)
	}

	user, err := getUser(0)
	if err != nil {
		t.Fatalf("%v", err)
	}

	require.Equal(t, user, models.User{
		ID:         1234,
		ScreenName: "NimartemX13",
		FirstName:  "Artem",
		LastName:   "Nizamov",
		Sex:        2,
		City: models.City{
			Title: "Nefteyugansk",
		},
		Followers: nil,
		Subscriptions: models.Subscriptions{
			Users: []models.User{
				{
					ID:         12345,
					ScreenName: "Puritanin",
					FirstName:  "Pavel",
					LastName:   "Demukhametov",
					Sex:        2,
					City: models.City{
						Title: "Nizhnevartovsk",
					},
				},
				{
					ID:         123456,
					ScreenName: "PostalDude",
					FirstName:  "Yakov",
					LastName:   "Zarembo",
					Sex:        2,
					City: models.City{
						Title: "Tyumen",
					},
				},
			},
			Groups: nil,
		},
	})
}

func TestDeleteUser(t *testing.T) {
	if err := deleteNode(1234); err != nil {
		t.Fatalf("%v", err)
	}
}

func getNodes() ([]models.Node, error) {
	resp, err := resty.New().
		SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime).
		AddRetryAfterErrorCondition().
		SetBaseURL(baseURL).
		R().Get("")

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.Status())
	}

	body := resp.Body()
	if body == nil {
		return nil, fmt.Errorf("body == nil")
	}

	var nodes []models.Node
	if err := sonic.Unmarshal(body, &nodes); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return nodes, nil
}

func getUser(id int64) (models.User, error) {
	resp, err := resty.New().
		SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime).
		AddRetryAfterErrorCondition().
		SetBaseURL(baseURL).
		R().SetPathParam("id", strconv.FormatInt(id, 10)).Get("{id}")

	if err != nil {
		return models.User{}, fmt.Errorf("%v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return models.User{}, fmt.Errorf("%v", resp.Status())
	}

	body := resp.Body()
	if body == nil {
		return models.User{}, fmt.Errorf("body == nil")
	}

	var user models.User
	if err := sonic.Unmarshal(body, &user); err != nil {
		return models.User{}, fmt.Errorf("%v", err)
	}

	return user, nil
}

func getGroup(id int64) (models.GroupWithSubscribers, error) {
	resp, err := resty.New().
		SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime).
		AddRetryAfterErrorCondition().
		SetBaseURL(baseURL).
		R().SetPathParam("id", strconv.FormatInt(id, 10)).Get("{id}")

	if err != nil {
		return models.GroupWithSubscribers{}, fmt.Errorf("%v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return models.GroupWithSubscribers{}, fmt.Errorf("%v", resp.Status())
	}

	body := resp.Body()
	if body == nil {
		return models.GroupWithSubscribers{}, fmt.Errorf("body == nil")
	}

	var group models.GroupWithSubscribers
	if err := sonic.Unmarshal(body, &group); err != nil {
		return models.GroupWithSubscribers{}, fmt.Errorf("%v", err)
	}

	return group, nil
}

func deleteNode(id int64) error {
	resp, err := resty.New().
		SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime).
		AddRetryAfterErrorCondition().
		SetAuthToken(token).
		SetBaseURL(baseURL).
		R().Delete("3")

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%v", resp.Status())
	}

	body := resp.Body()
	if body == nil {
		return fmt.Errorf("body == nil")
	}

	return nil
}
