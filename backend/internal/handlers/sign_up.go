package handlers

import (
	"errors"
	"net/http"
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
)

func (s *Server) SignUp(c *gin.Context) {
	var body api.SignUpJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Malformed JSON payload: " + err.Error(),
		})
		return
	}

	errorDetails := s.AuthService.ValidateSignUpCredentials(body)
	if len(errorDetails) > 0 {
		c.JSON(http.StatusUnprocessableEntity, api.ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Request validation failed.",
			Details: &errorDetails,
		})
		return
	}

	account, session, err := s.AuthService.SignUp(c, body)
	if err != nil {
		var apiErr *api_error.ApiError
		if errors.As(err, &apiErr) {
			api_error.HandleApiError(c, *apiErr)
			return
		}

		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Sign up failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, api.SignUpOkResponse{
		Account: api.Account{
			Id:          account.Id,
			Email:       account.Email,
			DisplayName: account.DisplayName,
		},
		Session: session,
	})
}
