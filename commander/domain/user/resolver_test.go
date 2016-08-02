package user_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/user"
)

func TestResolver(t *testing.T) {
	convey.Convey("Given a User Resolver", t, func() {
		userResolver := user.UserResolver{}

		convey.Convey("When it receives an emailAddress and password", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "foobar1234",
			})
			body := ioutil.NopCloser(bytes.NewReader(b))

			user := &user.User{}

			err := userResolver.FromRequest(user, body)

			convey.Convey("It should resolve correctly", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(user.EmailAddress, convey.ShouldEqual, "foo@bar.com")
				convey.So(user.Password, convey.ShouldEqual, "foobar1234")
			})
		})

		convey.Convey("When it doesn't receive an emailAddress and password", func() {
			b1, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
			})
			body1 := ioutil.NopCloser(bytes.NewReader(b1))
			user1 := &user.User{}
			err1 := userResolver.FromRequest(user1, body1)

			b2, _ := json.Marshal(map[string]string{
				"password": "foobar1234",
			})
			body2 := ioutil.NopCloser(bytes.NewReader(b2))
			user2 := &user.User{}
			err2 := userResolver.FromRequest(user2, body2)

			convey.Convey("It should return an error", func() {
				convey.So(err1, convey.ShouldNotBeNil)
				convey.So(user1.EmailAddress, convey.ShouldBeEmpty)
				convey.So(user1.Password, convey.ShouldBeEmpty)

				convey.So(err2, convey.ShouldNotBeNil)
				convey.So(user2.EmailAddress, convey.ShouldBeEmpty)
				convey.So(user2.Password, convey.ShouldBeEmpty)
			})
		})

		convey.Convey("When it receives the firstName and lastName as well", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "foobar1234",
				"firstName":    "foo",
				"lastName":     "bar",
			})
			body := ioutil.NopCloser(bytes.NewReader(b))

			user := &user.User{}

			err := userResolver.FromRequest(user, body)

			convey.Convey("It should resolve correctly", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(user.EmailAddress, convey.ShouldEqual, "foo@bar.com")
				convey.So(user.Password, convey.ShouldEqual, "foobar1234")
				convey.So(user.FirstName, convey.ShouldEqual, "foo")
				convey.So(user.LastName, convey.ShouldEqual, "bar")
			})
		})

		convey.Convey("When it receives malformed JSON", func() {
			b := []byte("malformedJSON")
			body := ioutil.NopCloser(bytes.NewReader(b))
			user := &user.User{}
			err := userResolver.FromRequest(user, body)

			convey.Convey("It should return an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(user.EmailAddress, convey.ShouldBeEmpty)
				convey.So(user.Password, convey.ShouldBeEmpty)
				convey.So(user.FirstName, convey.ShouldBeEmpty)
				convey.So(user.LastName, convey.ShouldBeEmpty)
			})
		})
	})
}
