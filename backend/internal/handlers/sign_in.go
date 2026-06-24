package handlers

import (
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

	account, session, apiErr := s.AuthService.SignIn(c, body)
	if apiErr != nil {
		api_error.HandleApiError(c, *apiErr)
		return
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
