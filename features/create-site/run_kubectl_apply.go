package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"tgs-automation/internal/log"

	"github.com/gin-gonic/gin"
)

// 主處理函數，重構以增強可讀性和可維護性
func runKubectlApply(lobby *LobbyInfo, request CreateSiteRequest) (int, map[string]any) {
	envMap := prepareKubectlEnv(lobby, request)
	templateContent, err := os.ReadFile(fmt.Sprintf("./deployment/lobby/lobby-%v.yaml", request.NameSpace))
	if err != nil {
		return logAndReturnError("Error reading template file", err)
	}

	config, err := applyKubectlTemplate(templateContent, envMap)
	if err != nil {
		return logAndReturnError("Error executing template", err)
	}

	fmt.Println(string(config))

	if err := os.WriteFile("target.yaml", config, 0644); err != nil {
		return logAndReturnError("Error writing file", err)
	}

	if err := executeKubectl("target.yaml"); err != nil {
		return logAndReturnError("Error running kubectl apply", err)
	}

	return http.StatusOK, gin.H{"status": "success"}
}

// 準備環境變量映射
func prepareKubectlEnv(lobby *LobbyInfo, request CreateSiteRequest) map[string]string {
	return map[string]string{
		"lobby":    fmt.Sprintf("lobby-%v-%v", strings.ToLower(request.BrandCode), strings.ToLower(request.LobbyTemplate)),
		"image":    lobby.DockerImage,
		"lang":     "en-US",
		"domain":   request.Domain,
		"token":    lobby.BrandToken,
		"currency": "CNY",
		"brand":    request.BrandCode,
	}
}

func logAndReturnError(message string, err error) (int, map[string]any) {
	log.LogError(fmt.Sprintf("%s: %v", message, err))
	return http.StatusInternalServerError, gin.H{"error": message, "details": err.Error()}
}

// preprocessTemplate 將 $var 變量替換為 {{.Var}}，以適配 Go 的模板語法
func preprocesskubectlTemplate(content []byte) string {
	// 將 $var 替換為 {{.Var}}
	// 注意這裡的替換邏輯可能需要根據實際情況調整以避免錯誤替換
	s := string(content)
	s = strings.ReplaceAll(s, "$lobby", "{{.lobby}}")
	s = strings.ReplaceAll(s, "$domain", "{{.domain}}")
	s = strings.ReplaceAll(s, "$image", "{{.image}}")
	s = strings.ReplaceAll(s, "$token", "{{.token}}")
	s = strings.ReplaceAll(s, "$currency", "{{.currency}}")
	s = strings.ReplaceAll(s, "$lang", "{{.lang}}")
	s = strings.ReplaceAll(s, "$brand", "{{.brand}}")
	return s
}

// 應用模板
func applyKubectlTemplate(templateContent []byte, envMap map[string]string) ([]byte, error) {
	modifiedTemplate := preprocesskubectlTemplate(templateContent)
	tmpl, err := template.New("config").Parse(modifiedTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, envMap); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 執行 kubectl
func executeKubectl(file string) error {
	cmd := exec.Command("kubectl", "apply", "-f", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
