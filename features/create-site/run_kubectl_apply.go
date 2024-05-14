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

	"github.com/gin-gonic/gin"
)

// 主處理函數，重構以增強可讀性和可維護性
func runKubectlApply(lobby *LobbyInfo, request CreateSiteRequest) (int, map[string]any) {
	envMap := prepareEnv(lobby, request)
	templateContent, err := os.ReadFile(fmt.Sprintf("lobby-%v.yaml", request.NameSpace))
	if err != nil {
		return logAndReturnError("Error reading template file", err)
	}

	config, err := applyTemplate(templateContent, envMap)
	if err != nil {
		return logAndReturnError("Error executing template", err)
	}

	if err := os.WriteFile("target.yaml", config, 0644); err != nil {
		return logAndReturnError("Error writing file", err)
	}

	if err := executeKubectl("target.yaml"); err != nil {
		return logAndReturnError("Error running kubectl apply", err)
	}

	return http.StatusOK, gin.H{"status": "success"}
}

// 準備環境變量映射
func prepareEnv(lobby *LobbyInfo, request CreateSiteRequest) map[string]string {
	return map[string]string{
		"lobby":  fmt.Sprintf("lobby-%v", request.BrandCode),
		"image":  lobby.DockerImage,
		"lang":   "en-US",
		"domain": request.Domain,
		"token":  lobby.BrandToken,
	}
}

func logAndReturnError(message string, err error) (int, map[string]any) {
	log.LogError(fmt.Sprintf("%s: %v", message, err))
	return http.StatusInternalServerError, gin.H{"error": message, "details": err.Error()}
}

// 應用模板
func applyTemplate(templateContent []byte, envMap map[string]string) ([]byte, error) {
	tmpl, err := template.New("config").Parse(string(templateContent))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, envMap); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 寫入到文件
func writeToFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

// 執行 kubectl
func executeKubectl(file string) error {
	cmd := exec.Command("kubectl", "apply", "-f", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
