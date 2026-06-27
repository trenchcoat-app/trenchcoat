package handlers

import (
	"net/http"
	"trenchcoat/internal/api"
	"trenchcoat/internal/dto/httperror"
	"trenchcoat/internal/services/cookie"

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

	signInResponse, httpErr := s.AuthService.SignIn(c, body)
	if httpErr != nil {
		httperror.HandleHttpError(c, httpErr)
		return
	}

	cookie.SetSessionCookie(c, *signInResponse.Session.SessionToken, *signInResponse.Session.ExpiresAt)

	// Write final response
	c.JSON(
		http.StatusOK,
		api.SignInOkResponse{
			Account: api.Account{
				Id:          signInResponse.Account.Id,
				Email:       body.Email,
				DisplayName: signInResponse.Account.DisplayName,
			},
		},
	)
}
