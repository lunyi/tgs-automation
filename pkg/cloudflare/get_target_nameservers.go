package cloudflare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"tgs-automation/internal/util"
)

func GetTargetNameServers(domain string) (string, error) {
	config := util.GetConfig()
	// Get Cloudflare name servers
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?per_page=1000&name=%s", domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.CloudflareToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Cloudflare data:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	nameServers := ""
	if res, ok := result["result"].([]interface{}); ok && len(res) > 0 {
		if ns, ok := res[0].(map[string]interface{})["name_servers"].([]interface{}); ok {
			var nsList []string
			for _, v := range ns {
				nsList = append(nsList, v.(string))
			}
			nameServers = strings.Join(nsList, " ")
		}
	}

	return nameServers, nil
}
