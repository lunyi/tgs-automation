package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"

	_ "tgs-automation/features/data-retrieve-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/swaggerFiles"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()
	router.GET("/healthz", healthCheckHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/token", tokenHandler)
	router.POST("/site", AuthMiddleware(), createLobbyHandler)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

// tokenHandler creates and sends a new JWT token
func tokenHandler(c *gin.Context) {
	tokenString, err := GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func createLobbyHandler(c *gin.Context) {
	config := util.GetConfig()
	request, err := parseRequestBody(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	httpCode, response := handleCreateLobbyInfo(request, config)

	if httpCode != http.StatusOK {
		c.JSON(httpCode, response)
		return
	}

	httpCode, response = handleKubectlApply(c)

	c.JSON(httpCode, response)
}

func parseRequestBody(c *gin.Context) (CreateSiteRequest, error) {
	var request CreateSiteRequest
	if err := c.BindJSON(&request); err != nil {
		log.LogError(fmt.Sprintf("Error parsing request body: %v", err))
		return request, err
	}
	log.LogInfo(fmt.Sprintf("Request: %v", request))
	return request, nil
}

func handleKubectlApply(lobby LobbyInfo, request CreateSiteRequest) (int, map[string]any) {
	templateFile := fmt.Sprintf("lobby-%v.yaml", request.NameSpace)

	content, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.LogError(fmt.Sprintf("Error reading template file:", err))
		return http.StatusInternalServerError, gin.H{"error": "Error reading template file", "details": err.Error()}
	}

	// 創建模板並應用環境變量
	tmpl, err := template.New("config").Parse(string(content))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		log.LogError(fmt.Sprintf("Error reading template file:", err))
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, os.Environ())
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// 將結果寫入新文件
	outputFile := "target.yaml"
	err = ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	// 執行 kubectl apply 命令
	cmd := exec.Command("kubectl", "apply", "-f", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running kubectl apply:", err)
	}
}

func handleCreateLobbyInfo(request CreateSiteRequest, config util.TgsConfig) (int, map[string]any) {
	// Fetch Docker image
	dockerhubService := NewDockerImageService(config.Dockerhub)
	image, err := dockerhubService.FetchDockerImage(request.LobbyTemplate)

	if err != nil {
		log.LogError(fmt.Sprintf("Error fetching docker image %v", err))
		return http.StatusInternalServerError, gin.H{"error": "Could not get docker image", "details": err.Error()}
	}
	fmt.Println("Image:", image)

	// Get brand ID
	brandId, err := postgresql.GetBrandId(config.CreateSiteDb, "MOPH")
	if err != nil {
		log.LogError(fmt.Sprintf("Error getting brand id %v", err))
		return http.StatusInternalServerError, gin.H{"error": "Could not get brand id", "details": err.Error()}
	}
	fmt.Println("Brand ID:", brandId)

	// Get brand token
	token, err := getBrandToken(brandId, "staging", config.ApiUrl.BrandCert)
	if err != nil {
		log.LogError(fmt.Sprintf("Error getting brand token %v", err))
		return http.StatusInternalServerError, gin.H{"error": "Could not get brand token", "details": err.Error()}
	}

	lobby := &LobbyInfo{
		BrandToken:  token,
		DockerImage: image}
	return http.StatusOK, gin.H{"lobby": lobby}
}

type CreateSiteRequest struct {
	BrandCode     string `json:"brandCode"`
	LobbyTemplate string `json:"lobbyTemplate"`
	Domain        string `json:"domain"`
	NameSpace     string `json:"namespace"`
}

type LobbyInfo struct {
	BrandToken  string
	DockerImage string
}
