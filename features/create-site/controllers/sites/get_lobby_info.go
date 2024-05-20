package sites

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"

	"github.com/gin-gonic/gin"
)

type LobbyInfo struct {
	BrandToken  string
	DockerImage string
}

func GetLobbyInfo(request CreateSiteRequest, config util.TgsConfig) (int, map[string]any) {
	// Fetch Docker image
	dockerhubService := NewDockerImageService(config.Dockerhub)
	image, err := dockerhubService.FetchDockerImage(request.LobbyTemplate)

	if err != nil {
		log.LogError(fmt.Sprintf("Error fetching docker image %v", err))
		return http.StatusInternalServerError, gin.H{"error": "Could not get docker image", "details": err.Error()}
	}
	fmt.Println("Image:", image)

	// Get brand ID
	brandId, err := postgresql.GetBrandId(config.CreateSiteDb, request.BrandCode)
	if err != nil {
		log.LogError(fmt.Sprintf("Error getting brand id %v", err))
		return http.StatusInternalServerError, gin.H{"error": "Could not get brand id", "details": err.Error()}
	}
	fmt.Println("Brand ID:", brandId)

	// Get brand token
	token, err := GetBrandToken(brandId, request.NameSpace, config.ApiUrl.BrandCert)
	if err != nil {
		log.LogError(fmt.Sprintf("Error getting brand token %v", err))
		return http.StatusInternalServerError, gin.H{"error": "Could not get brand token", "details": err.Error()}
	}

	lobby := &LobbyInfo{
		BrandToken:  token,
		DockerImage: image}
	return http.StatusOK, gin.H{"lobby": lobby}
}

func ConvertToLobbyInfo(response map[string]any) (*LobbyInfo, error) {
	// 將 map 轉換為 JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("error marshalling map to JSON: %v", err)
	}

	// 定義一個中間結果的結構
	var result struct {
		Lobby LobbyInfo `json:"lobby"`
	}

	// 從 JSON 中解析回中間結果的結構
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON to struct: %v", err)
	}

	// 返回 LobbyInfo 結構體
	return &result.Lobby, nil
}
