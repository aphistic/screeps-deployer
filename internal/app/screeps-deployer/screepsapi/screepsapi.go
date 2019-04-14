package screepsapi

import (
	"fmt"
	"net/http"
)

type basicResponse struct {
	OK int `json:"ok"`
}

type ClientOption func(*Client)

func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.token = token
	}
}

type Client struct {
	client *http.Client

	token string

	scheme string
	host   string
}

func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		client: &http.Client{},

		scheme: "https",
		host:   "screeps.com",
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) baseURL() string {
	return fmt.Sprintf("%s://%s/", c.scheme, c.host)
}
