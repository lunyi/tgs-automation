package sites

import (
	"fmt"
	"net/http"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateSiteHandler(c *gin.Context) {
	config := util.GetConfig()
	request, err := parseRequestBody(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	httpCode, response := GetLobbyInfo(request, config)

	if httpCode != http.StatusOK {
		c.JSON(httpCode, response["lobby"])
		return
	}

	lobbyInfo, err := ConvertToLobbyInfo(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not convert lobby info", "details": err.Error()})
		return
	}

	httpCode, response = RunKubectlApply(lobbyInfo, request)
	if httpCode != http.StatusOK {
		c.JSON(httpCode, response)
		return
	}

	c.JSON(httpCode, response)
}

type CreateSiteRequest struct {
	BrandCode     string `json:"brandCode" validate:"required"`
	LobbyTemplate string `json:"lobbyTemplate" validate:"required"`
	Domain        string `json:"domain" validate:"required"`
	NameSpace     string `json:"namespace" validate:"required,oneof=dev staging prod"`
}

// 定义一个全局的验证器
var validate = validator.New()

func parseRequestBody(c *gin.Context) (CreateSiteRequest, error) {
	var request CreateSiteRequest
	if err := c.BindJSON(&request); err != nil {
		log.LogError(fmt.Sprintf("Error parsing request body: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return request, err
	}

	// 使用驗證器驗證請求體
	if err := validate.Struct(&request); err != nil {
		validationErrors := getValidationErrors(err)
		errMsg := fmt.Sprintf("Validation errors: %v", validationErrors)
		log.LogError(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return request, fmt.Errorf(errMsg)
	}

	log.LogInfo(fmt.Sprintf("Request: %v", request))
	return request, nil
}

// getValidationErrors 返回驗證錯誤信息
func getValidationErrors(err error) []string {
	var validationErrors []string
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, fmt.Sprintf("Field %s failed validation: %s", err.StructField(), err.Tag()))
	}
	return validationErrors
}
