package user_test

import (
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/user"
)

func TestProvider(t *testing.T) {
	convey.Convey("Given a User Provider", t, func() {
		userProvider := user.UserProvider{}

		convey.Convey("It should create a new User", func() {
			res := userProvider.New()

			convey.So(reflect.TypeOf(res).String(), convey.ShouldEqual, "*user.User")
		})
	})
}
