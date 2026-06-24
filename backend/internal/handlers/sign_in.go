package handlers

import (
	"errors"
	"net/http"
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
)

func (s *Server) SignIn(c *gin.Context) {
	var body api.SignInJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Malformed JSON payload: " + err.Error(),
		})
		return
	}

	errorDetails := s.AuthService.ValidateSignInCredentials(body)
	if len(errorDetails) > 0 {
		c.JSON(http.StatusUnprocessableEntity, api.ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Request validation failed.",
			Details: &errorDetails,
		})
		return
	}

	account, session, err := s.AuthService.SignIn(c, body)

	if err != nil {
		var apiErr *api_error.ApiError
		if errors.As(err, &apiErr) {
			api_error.HandleApiError(c, *apiErr)
			return
		}

		// Fallback
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create session: " + err.Error(),
		})
	}

	// Write final response
	c.JSON(
		http.StatusOK,
		api.SignInOkResponse{
			Account: api.Account{
				Id:          account.Id,
				Email:       body.Email,
				DisplayName: account.DisplayName,
			},
			Session: api.Session{
				Token:     session.Token,
				ExpiresAt: session.ExpiresAt,
				AccountId: account.Id,
			},
		},
	)
}
