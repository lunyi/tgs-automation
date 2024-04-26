package dockerhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DockerTagsResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

// GetLatestTagByImage fetches the latest tag for a given image from Docker Hub.
func (dhc *DockerHubClient) GetLatestTagByImage(image string) (string, error) {
	url := dhc.buildTagURL(image)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+dhc.Token)

	resp, err := dhc.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	tagName, err := dhc.extractTagName(resp)
	if err != nil {
		return "", err
	}

	return tagName, nil
}

// buildTagURL constructs the URL to fetch the latest tag of an image.
func (dhc *DockerHubClient) buildTagURL(image string) string {
	return fmt.Sprintf("%s/repositories/%s/%s/tags?page=1&page_size=1", dhc.Config.BaseUrl, dhc.Config.Username, image)
}

// extractTagName reads the response body and extracts the tag name.
func (dhc *DockerHubClient) extractTagName(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	var tagsResp DockerTagsResponse
	if err := json.Unmarshal(body, &tagsResp); err != nil {
		return "", fmt.Errorf("unmarshaling response: %w", err)
	}

	if len(tagsResp.Results) == 0 {
		return "", fmt.Errorf("no tags found for image %s", dhc.Config.Username)
	}
	return tagsResp.Results[0].Name, nil
}
