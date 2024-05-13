package main

import (
	"fmt"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
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

	brandId, err := postgresql.GetBrandId(config.Postgresql, "MOPH")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Brand ID:", brandId)
}
