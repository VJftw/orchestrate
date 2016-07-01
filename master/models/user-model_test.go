package models

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUserModel(t *testing.T) {
	convey.Convey("Given a User", t, func() {
		user := User{
			HashID:       "abcdef",
			EmailAddress: "foo@bar.com",
			Password:     "abc1234",
			FirstName:    "Foo",
			LastName:     "Bar",
		}
		user.ID = 2

		convey.Convey("It should return a serializable map", func() {
			convey.So(user.ToMap()["id"], convey.ShouldEqual, "abcdef")
			convey.So(user.ToMap()["emailAddress"], convey.ShouldEqual, "foo@bar.com")
			convey.So(user.ToMap()["firstName"], convey.ShouldEqual, "Foo")
			convey.So(user.ToMap()["lastName"], convey.ShouldEqual, "Bar")
		})

		convey.Convey("It should return the ID", func() {
			convey.So(user.GetID(), convey.ShouldEqual, 2)
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
