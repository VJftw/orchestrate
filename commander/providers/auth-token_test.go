package providers

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/models"
)

func TestAuthTokenProvider(t *testing.T) {
	convey.Convey("Given an Auth Token Provider", t, func() {
		authTokenProvider := NewAuthToken()

		convey.Convey("When I call NewFromUser() with a User", func() {
			user := models.User{
				UUID: "1123-ads12-12312asd-as",
			}
			authToken, err := authTokenProvider.NewFromUser(&user)

			convey.Convey("Then an AuthToken should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(authToken.Token, convey.ShouldNotBeEmpty)
			})
		})
	})
}
