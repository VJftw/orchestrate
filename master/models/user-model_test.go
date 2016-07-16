package models

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUserModel(t *testing.T) {
	convey.Convey("Given a User", t, func() {
		user := User{
			UUID:         "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			EmailAddress: "foo@bar.com",
			Password:     "abc1234",
			FirstName:    "Foo",
			LastName:     "Bar",
		}

		convey.Convey("It should return a serializable map", func() {
			convey.So(user.ToMap()["uuid"], convey.ShouldEqual, "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
			convey.So(user.ToMap()["emailAddress"], convey.ShouldEqual, "foo@bar.com")
			convey.So(user.ToMap()["firstName"], convey.ShouldEqual, "Foo")
			convey.So(user.ToMap()["lastName"], convey.ShouldEqual, "Bar")
		})

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
