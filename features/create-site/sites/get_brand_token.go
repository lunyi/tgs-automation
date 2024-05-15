package sites

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"tgs-automation/internal/log"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func GetBrandToken(brandId, namespace, brandCertUrl string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//url := fmt.Sprintf("http://brand-cert-api.%v/cert/%v", namespace, brandId)
	url := fmt.Sprintf("%v%v", brandCertUrl, brandId)

	log.LogInfo(fmt.Sprintf("Requesting brand token from %v", url))

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("創建請求失敗: %v", err)
	}

	req.Header.Add("Accept", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("發送請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("非預期的響應狀態碼: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("讀取響應內容失敗: %v", err)
	}

	fmt.Println("Status:", resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return string(body), nil
	}

	return "", fmt.Errorf("非預期的響應狀態碼: %v", resp.Status)
}
