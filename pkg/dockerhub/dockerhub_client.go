package dockerhub

import (
	"net/http"
	"tgs-automation/internal/util"
)

// DockerHubClient handles interactions with the Docker Hub API.
type DockerHubClient struct {
	Config util.DockerhubConfig
	Client *http.Client
	Token  string
}

type DockerHubClientInterface interface {
	GetLatestTagByImage(image string) (string, error)
}

// NewDockerHubClient creates a new instance of DockerHubClient with the provided configuration.
func NewDockerHubClient(config util.DockerhubConfig) *DockerHubClient {
	token, err := Login(config)
	if err != nil {
		panic(err)
	}
	return &DockerHubClient{
		Token:  token,
		Config: config,
		Client: &http.Client{}, // This could be customized for things like timeouts, retries, etc.
	}
}
