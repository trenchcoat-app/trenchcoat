package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"trenchcoat/internal/api"
	"trenchcoat/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestPingRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := server.GetRouter(nil)

	httpRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	router.ServeHTTP(httpRecorder, req)

	var pingResponse api.PingResponse
	err := json.Unmarshal(httpRecorder.Body.Bytes(), &pingResponse)
	require.NoError(t, err)

	assert.Equal(t, 200, httpRecorder.Code)
	assert.Equal(t, api.PingResponse{Message: "pong"}, pingResponse)
}
