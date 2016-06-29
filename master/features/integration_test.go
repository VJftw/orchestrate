package integration

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/features/user-management"
)

func TestIntegration(t *testing.T) {
	Convey("Given some integer with a starting value", t, func() {
		x := 1

		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})

	userManagement := integrationUserManagement.UserManagementTests{}

	userManagement.UserRegistration(t)
}
