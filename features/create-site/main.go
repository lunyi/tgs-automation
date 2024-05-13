package main

import (
	"fmt"
	"tgs-automation/internal/util"
)

func main() {

	// Get the config
	config := util.GetConfig()
	dockerhubService := NewDockerImageService(config.Dockerhub)
	image, err := dockerhubService.FetchDockerImage("T23")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Image:", image)
}
