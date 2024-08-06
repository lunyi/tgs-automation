package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockNatsPublish "tgs-automation/internal/util/mocks"
	mockNamecheap "tgs-automation/pkg/namecheap/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateDomainHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNamecheapApi := mockNamecheap.NewMockNamecheapApi(ctrl)
	mockNats := mockNatsPublish.NewMockNatsPublisherService(ctrl)

	mockNats.EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		Return(nil)

	mockNamecheapApi.EXPECT().
		GetCouponCode(gomock.Any()).
		Return("PROMOCODE", nil)

	mockNamecheapApi.EXPECT().
		CreateDomain(gomock.Any(), "example.com", "PROMOCODE").
		Return(`<ApiResponse Status="OK"><CommandResponse><DomainCreateResult ChargedAmount="8.88"/></CommandResponse></ApiResponse>`, nil)

	mockNamecheapApi.EXPECT().
		GetBalance(gomock.Any()).
		Return(`<ApiResponse><CommandResponse><UserGetBalancesResult Currency="USD" AvailableBalance="100.00"/></CommandResponse></ApiResponse>`, nil)

	// Setup the router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/domain", CreateDomainHandler(mockNamecheapApi, mockNats))

	// Create a request body
	requestBody, _ := json.Marshal(map[string]interface{}{
		"domain": "example.com",
		"chatid": "12345",
	})

	// Create a HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/domain", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a HTTP recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response["status"])
	assert.NotEmpty(t, response["message"])
}
