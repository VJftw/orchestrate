package user

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/features/utils"
)

type UserAuthTests struct {
}

func (uAT UserAuthTests) UserAuthentication(t *testing.T, apiClient utils.APIClient) {

	convey.Convey("User Authentication tests", func() {
		convey.Convey("When I send a POST request to /v1/auth with valid User credentials", func() {
			body := map[string]interface{}{
				"emailAddress": "foo@bar.com",
				"password":     "abcd1234",
			}

			apiClient.RequestWithBody("POST", "/v1/auth", body)

			convey.Convey("Then the Response should contain a JWT", func() {
				convey.So(apiClient.ResponseStatus, convey.ShouldEqual, 201)
				var authResp struct {
					AuthToken string `json:"authToken"`
				}
				apiClient.UnmarshalTo(&authResp)

				convey.So(authResp.AuthToken, convey.ShouldNotBeEmpty)
			})
		})

	})

}
