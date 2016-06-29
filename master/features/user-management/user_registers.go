package integrationUserManagement

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type UserManagementTests struct {
}

func (uM *UserManagementTests) UserRegistration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	Convey("When I Send a POST request to /v1/users with valid User data", t, func() {

		Convey("Then the User should be created", func() {

		})
	})
}
