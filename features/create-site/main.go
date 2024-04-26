package main

import (
	"fmt"
	"strings"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/dockerhub"
)

func main() {

	// Get the config
	config := util.GetConfig()

	fmt.Println("Dockerhub config:", config.Dockerhub)

	client := dockerhub.NewDockerHubClient(config.Dockerhub)

	image := getDockerImageByLobbyTemplate("t1")

	tag, err := client.GetLatestTagByImage(image)
	if err != nil {
		fmt.Println("Failed to fetch the latest tag for image:", err)
		return
	}

	image = fmt.Sprintf("%s/%s:%s", config.Dockerhub.Username, image, tag)
	fmt.Println("Image:", image)
}

func getDockerImageByLobbyTemplate(lobbyTemplate string) string {
	lobbyTemplate = strings.ToLower(lobbyTemplate)
	if lobbyTemplate <= "t21" {
		return fmt.Sprintf("lobby-%s", lobbyTemplate)
	}
	return fmt.Sprintf("tgs-web-%s", lobbyTemplate)
}
