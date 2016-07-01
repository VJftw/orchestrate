package user

import (
	"testing"

	"github.com/vjftw/orchestrate/master/features/utils"
)

type UserManagementTests struct {
}

func (uMT UserManagementTests) RunTests(t *testing.T, apiClient utils.APIClient) {
	registration := UserRegistrationTests{}
	registration.UserRegistration(t, apiClient)
	// auth := UserAuthTests{}
	// auth.UserAuthentication(t, apiClient)

}
