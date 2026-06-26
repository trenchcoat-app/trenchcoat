package api_error_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestHandleApiError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	apiErr := api_error.ApiError{
		Status:  http.StatusConflict,
		Code:    "EMAIL_ALREADY_EXISTS",
		Message: "An account with this email address already exists.",
	}

	api_error.HandleApiError(c, apiErr)

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", resp.Code)
	assert.Equal(t, "An account with this email address already exists.", resp.Message)
}
