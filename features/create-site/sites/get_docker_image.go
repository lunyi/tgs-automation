package sites

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/dockerhub"
)

type DockerImageFetcher interface {
	FetchDockerImage(template string) (string, error)
}

type DockerImageService struct {
	config util.DockerhubConfig
	client *dockerhub.DockerHubClient
}

func NewDockerImageService(config util.DockerhubConfig) DockerImageFetcher {
	client := dockerhub.NewDockerHubClient(config)
	return &DockerImageService{config: config, client: client}
}

func (service *DockerImageService) FetchDockerImage(template string) (string, error) {
	log.LogInfo(fmt.Sprintf("Fetching docker image for template: %s", template))
	image, err := getDockerImageByLobbyTemplate(template)
	if err != nil {
		return "", fmt.Errorf("invalid template format: %w", err)
	}

	log.LogInfo(fmt.Sprintf("Fetching the latest tag for image: %s", image))

	tag, err := service.client.GetLatestTagByImage(image)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the latest tag for image: %w", err)
	}

	image = fmt.Sprintf("%s/%s:%s", service.config.Username, image, tag)
	return image, nil
}

func getDockerImageByLobbyTemplate(lobbyTemplate string) (string, error) {
	templateNumber, err := getTemplateNumber(lobbyTemplate)

	if err != nil {
		return "", fmt.Errorf("invalid template number: %w", err)
	}

	lobbyTemplate = strings.ToLower(lobbyTemplate)
	if templateNumber <= 21 {
		return fmt.Sprintf("lobby-%s", lobbyTemplate), nil
	}
	return fmt.Sprintf("tgs-web-%s", strings.ToLower(lobbyTemplate)), nil
}

func getTemplateNumber(template string) (int, error) {
	log.LogInfo(fmt.Sprintf("Checking the template number: %s", template))
	regexPattern := `^[tT][0-9]`
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return -1, fmt.Errorf("error compiling regex: %v", err)
	}
	num, err := strconv.Atoi(template[1:])

	if err != nil {
		return -1, fmt.Errorf("error getting a number: %v", err)
	}
	log.LogInfo(fmt.Sprintf("checking regex: %v", regex))
	return num, nil
}
