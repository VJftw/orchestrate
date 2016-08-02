package user_test

import (
	"github.com/vjftw/orchestrate/commander/domain/user"

	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUser(t *testing.T) {
	convey.Convey("Given a User Model", t, func() {
		user := user.User{
			UUID:         "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			EmailAddress: "foo@bar.com",
			Password:     "abc1234",
			FirstName:    "Foo",
			LastName:     "Bar",
		}

		convey.Convey("It should return the UUID", func() {
			convey.So(string(user.GetUUID()), convey.ShouldEqual, "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		})

		convey.Convey("It should encrypt the password", func() {
			convey.So(user.PasswordHash, convey.ShouldBeNil)
			user.EncryptPassword()
			convey.So(user.PasswordHash, convey.ShouldNotBeNil)

			convey.Convey("It should verify a set password", func() {
				user.Password = "wrongPassword"
				convey.So(user.VerifyPassword(), convey.ShouldBeFalse)

				user.Password = "abc1234"
				convey.So(user.VerifyPassword(), convey.ShouldBeTrue)
			})
		})

	})
}
