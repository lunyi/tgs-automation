package sites

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"text/template"
	k8s "tgs-automation/internal/kubernetes"
	"tgs-automation/internal/log"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

// 主處理函數，重構以增強可讀性和可維護性
func RunKubectlApply(lobby *LobbyInfo, request CreateSiteRequest) (int, map[string]any) {
	envMap := prepareKubectlEnv(lobby, request)
	templateContent, err := os.ReadFile(fmt.Sprintf("./deployment/lobby/lobby-%v.yaml", request.NameSpace))
	if err != nil {
		return LogAndReturnApiError("Error reading template file", err)
	}

	modifiedTemplateBytes, err := getModifiedTemplate(templateContent, envMap)
	if err != nil {
		return LogAndReturnApiError("Error executing template", err)
	}

	clientset, err := k8s.InitKubeClient(request.NameSpace)
	if err != nil {
		return LogAndReturnApiError("Error initializing Kubernetes client", err)
	}

	if err := applyYamlFile(clientset, string(modifiedTemplateBytes)); err != nil {
		return LogAndReturnApiError("Error applying YAML file", err)
	}

	return http.StatusOK, gin.H{"status": "success"}
}

func applyYamlFile(clientset *kubernetes.Clientset, templateContent string) error {
	docs := strings.Split(templateContent, "---")
	for _, doc := range docs {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		if err := k8s.ApplyYamlDocument(clientset, doc); err != nil {
			return err
		}
	}

	return nil
}

// 準備環境變量映射
func prepareKubectlEnv(lobby *LobbyInfo, request CreateSiteRequest) map[string]string {
	return map[string]string{
		"lobby":      fmt.Sprintf("lobby-%v-%v", strings.ToLower(request.BrandCode), strings.ToLower(request.LobbyTemplate)),
		"image":      lobby.DockerImage,
		"lang":       "en-US",
		"domain":     request.Domain,
		"token":      lobby.BrandToken,
		"currency":   "CNY",
		"brand":      request.BrandCode,
		"theme":      request.LobbyTemplate[1:],
		"wapemplate": fmt.Sprintf("wap%s", strings.ToUpper(request.LobbyTemplate)),
	}
}

func LogAndReturnApiError(message string, err error) (int, map[string]any) {
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
	s = strings.ReplaceAll(s, "$theme", "{{.theme}}")
	s = strings.ReplaceAll(s, "$waptemplate", "{{.waptemplate}}")
	return s
}

// 應用模板
func getModifiedTemplate(templateContent []byte, envMap map[string]string) ([]byte, error) {
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
