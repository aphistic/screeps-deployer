package screepsapi

import (
	"fmt"
	"net/http"
)

type basicResponse struct {
	OK int `json:"ok"`
}

type Client struct {
	client *http.Client

	scheme string
	host   string
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{},

		scheme: "https",
		host:   "screeps.com",
	}
}

func (c *Client) baseURL() string {
	return fmt.Sprintf("%s://%s", c.scheme, c.host)
}
