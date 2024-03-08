package namecheap

import (
	"bytes"
	"cdnetwork/internal/util"
	"cdnetwork/test/mocks"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {

	mockHttpClient := createMockHTTPClient(t)

	config := util.NamecheapConfig{
		NamecheapApiKey:   "NamecheapApiKey",
		NamecheapUsername: "NamecheapUsername",
		NamecheapPassword: "NamecheapPassword",
		NamecheapClientIp: "NamecheapClientIp",
	}
	api := New(config, mockHttpClient)
	namecheapClient, ok := api.(*NamecheapClient)
	if !ok {
		t.Error("New didn't return a NamecheapClient instance")
	}
	if namecheapClient.Config.NamecheapApiKey != config.NamecheapApiKey ||
		namecheapClient.Config.NamecheapUsername != config.NamecheapUsername ||
		namecheapClient.Config.NamecheapPassword != config.NamecheapPassword ||
		namecheapClient.Config.NamecheapClientIp != config.NamecheapClientIp {
		t.Error("Config mismatch")
	}
}

func TestNew_Assert(t *testing.T) {
	mockHttpClient := createMockHTTPClient(t)
	config := util.NamecheapConfig{
		NamecheapApiKey:   "NamecheapApiKey",
		NamecheapUsername: "NamecheapUsername",
		NamecheapPassword: "NamecheapPassword",
		NamecheapClientIp: "NamecheapClientIp",
	}
	api := New(config, mockHttpClient)
	namecheapClient, ok := api.(*NamecheapClient)
	assert.True(t, ok)
	assert.Equal(t, config.NamecheapApiKey, namecheapClient.Config.NamecheapApiKey)
	assert.Equal(t, config.NamecheapClientIp, namecheapClient.Config.NamecheapClientIp)
	assert.Equal(t, config.NamecheapPassword, namecheapClient.Config.NamecheapPassword)
	assert.Equal(t, config.NamecheapUsername, namecheapClient.Config.NamecheapUsername)
}

// A helper function to create a mock HTTP server response
// createMockResponse helps in creating a mock HTTP response
func createMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func createMockHTTPClient(t *testing.T) *mocks.MockHTTPClient {
	nextMonth := time.Now().AddDate(0, 1, 0)
	mockApiResponse := ApiResponse{
		Status: "OK",
		CommandResponse: CommandResponse{
			DomainGetListResult: DomainGetListResult{
				Domains: []Domain{
					{
						Name:    "expiredexample.com",
						Created: time.Now().Add(-96 * time.Hour).Format(dateFormat),
						Expires: time.Date(nextMonth.Year(), nextMonth.Month(), 2, 0, 0, 0, 0, nextMonth.Location()).Format(dateFormat),
					},
					{
						Name:    "futureexample.com",
						Created: time.Now().Add(-96 * time.Hour).Format(dateFormat),
						Expires: time.Date(nextMonth.Year(), nextMonth.Month(), 2, 0, 0, 0, 0, nextMonth.Location()).Format(dateFormat),
					},
				},
			},
			Paging: Paging{
				TotalItems:  2,
				CurrentPage: 1,
				PageSize:    2,
			},
		},
	}

	responseBody, err := xml.Marshal(mockApiResponse)
	if err != nil {
		t.Fatalf("Unable to marshal XML2: %v", err)
	}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("Get", mock.Anything).Return(createMockResponse(200, string(responseBody)), nil)

	return mockClient
}

func TestGetExpiredDomains(t *testing.T) {
	mockHttpClient := createMockHTTPClient(t)

	client := New(util.NamecheapConfig{
		NamecheapUsername: "user",
		NamecheapApiKey:   "key",
		NamecheapClientIp: "ip",
	}, mockHttpClient)

	// Mocking time.Now() to return a specific date if necessary
	// ...

	domains, err := client.GetExpiredDomains()

	// Check if an error occurred
	assert.NoError(t, err, "Unexpected error: %v", err)

	// Check if the result is as expected
	expectedDomain := FilteredDomain{Name: "example.com", Created: "01/01/2022", Expires: "01/02/2023"}
	assert.Equal(t, []FilteredDomain{expectedDomain}, domains, "Unexpected result")
}
