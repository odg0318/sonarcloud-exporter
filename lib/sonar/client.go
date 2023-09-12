package sonar

import (
	"fmt"
	"net/http"
	"time"
)

// ListOptions provide general API request options.
type ListOptions struct {
	Page    int
	PerPage int
}

// Client holds the variables we want to use for the connection.
type Client struct {
	ListOptions
	sonarConnectionString string
	client                *http.Client
	organization          string
	metrics               string
}

// NewClient creates a new SonarCloud client
func NewClient(token string, org string, metrics string) *Client {
	uri := fmt.Sprintf("https://%s@sonarcloud.io/api", token)

	return &Client{
		sonarConnectionString: uri,
		client:                &http.Client{Timeout: time.Second * 10},
		organization:          org,
		metrics:               metrics,
	}
}
