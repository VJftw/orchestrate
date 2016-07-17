package resolvers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/models"
)

func TestUserResolver(t *testing.T) {
	convey.Convey("Given a User Resolver", t, func() {
		userResolver := UserResolver{}

		convey.Convey("When the JSON is valid", func() {

			convey.Convey("And it only has Email Address and Password", func() {
				b, _ := json.Marshal(map[string]string{
					"emailAddress": "foo@bar.com",
					"password":     "foobar1234",
				})

				var user models.User

				err := userResolver.FromRequest(&user, ioutil.NopCloser(bytes.NewReader(b)))
				convey.Convey("Then it should not return an error", func() {
					convey.So(err, convey.ShouldBeNil)
				})
				convey.Convey("And the values should be set", func() {
					convey.So(user.EmailAddress, convey.ShouldEqual, "foo@bar.com")
					convey.So(user.Password, convey.ShouldEqual, "foobar1234")
				})
			})

			convey.Convey("Or it has all of the details", func() {
				b, _ := json.Marshal(map[string]string{
					"emailAddress": "foo@bar.com",
					"password":     "foobar1234",
					"firstName":    "Foo",
					"lastName":     "Bar",
				})

				var user models.User

				err := userResolver.FromRequest(&user, ioutil.NopCloser(bytes.NewReader(b)))
				convey.Convey("Then it should not return an error", func() {
					convey.So(err, convey.ShouldBeNil)
				})
				convey.Convey("And the values should be set", func() {
					convey.So(user.EmailAddress, convey.ShouldEqual, "foo@bar.com")
					convey.So(user.Password, convey.ShouldEqual, "foobar1234")
					convey.So(user.FirstName, convey.ShouldEqual, "Foo")
					convey.So(user.LastName, convey.ShouldEqual, "Bar")
				})
			})
		})
		convey.Convey("When the JSON is malformed", func() {
			b := []byte("malformedJSON")

			var user models.User
			err := userResolver.FromRequest(&user, ioutil.NopCloser(bytes.NewReader(b)))

			convey.Convey("Then it should return an Error", func() {
				convey.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
