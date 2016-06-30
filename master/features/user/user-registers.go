package user

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/features/utils"
)

type UserRegistrationTests struct {
}

func (uM UserRegistrationTests) UserRegistration(t *testing.T, apiClient utils.APIClient) {

	convey.Convey("User Registration tests", func() {
		convey.Convey("When I Send a POST request to /v1/users with valid User data", func() {
			body := map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "abcd1234",
			}
			err := apiClient.Post("/v1/users", body)
			if err != nil {
				t.Error(err)
			}
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Then the User should be created", func() {
				convey.So(apiClient.ResponseStatus, convey.ShouldEqual, 201)

				var userResp struct {
					UUID         string `json:"id"`
					EmailAddress string `json:"emailAddress"`
					FirstName    string `json:"firstName"`
					LastName     string `json:"lastName"`
				}
				apiClient.UnmarshalTo(&userResp)

				convey.So(userResp.UUID, convey.ShouldNotBeEmpty)
				convey.So(userResp.EmailAddress, convey.ShouldEqual, "foo@bar.com")
				convey.So(userResp.FirstName, convey.ShouldEqual, "")
				convey.So(userResp.LastName, convey.ShouldEqual, "")
			})
		})
	})

}
