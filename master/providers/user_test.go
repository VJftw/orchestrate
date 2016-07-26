package providers

import (
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUserProvider(t *testing.T) {
	convey.Convey("Given a User Provider", t, func() {
		userProvider := NewUser()

		convey.Convey("When I call New()", func() {
			res := userProvider.New()

			convey.Convey("Then it should return a new User", func() {
				convey.So(reflect.TypeOf(res).String(), convey.ShouldEqual, "*models.User")
				convey.So(res.UUID, convey.ShouldBeEmpty)
				convey.So(res.EmailAddress, convey.ShouldBeEmpty)
				convey.So(res.Password, convey.ShouldBeEmpty)
				convey.So(res.PasswordHash, convey.ShouldBeEmpty)
				convey.So(res.FirstName, convey.ShouldBeEmpty)
				convey.So(res.LastName, convey.ShouldBeEmpty)
			})
		})
	})
}
