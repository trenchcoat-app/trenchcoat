package testutil

import (
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func NewEmail() openapi_types.Email {
	return openapi_types.Email("e2e-test-" + uuid.NewString() + "@example.com")
}
