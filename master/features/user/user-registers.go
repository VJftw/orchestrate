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
		// var authUserHashID string

		convey.Convey("When I send a POST request to /v1/users with valid User data", func() {
			body := map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "abcd1234",
				"firstName":    "Foo",
				"lastName":     "Bar",
			}
			apiClient.RequestWithBody("POST", "/v1/users", body)

			convey.Convey("Then the User should be created", func() {
				convey.So(apiClient.ResponseStatus, convey.ShouldEqual, 201)

				var userResp struct {
					HashID       string `json:"id"`
					EmailAddress string `json:"emailAddress"`
					FirstName    string `json:"firstName"`
					LastName     string `json:"lastName"`
				}
				apiClient.UnmarshalTo(&userResp)

				convey.So(userResp.HashID, convey.ShouldNotBeEmpty)
				// authUserHashID = userResp.HashID
				convey.So(userResp.EmailAddress, convey.ShouldEqual, "foo@bar.com")
				convey.So(userResp.FirstName, convey.ShouldEqual, "Foo")
				convey.So(userResp.LastName, convey.ShouldEqual, "Bar")
			})
		})

		// convey.Convey("Given I am authenticated as foo@bar.com", func() {
		// 	body := map[string]string{
		// 		"emailAddress": "foo@bar.com",
		// 		"password":     "abcd1234",
		// 	}
		// 	apiClient.RequestWithBody("PUT", "/v1/auth", body)
		// 	var authResp struct {
		// 		AuthToken string `json:"authToken"`
		// 	}
		// 	apiClient.UnmarshalTo(authResp)
		//
		// 	apiClient.BearerToken = authResp.AuthToken
		//
		// 	convey.Convey(fmt.Sprintf("When I send a PUT request to /v1/users/%v with new valid user data", authUserHashID), func() {
		// 		body := map[string]string{
		// 			"firstName": "Foo2",
		// 			"lastName":  "Bar2",
		// 		}
		// 		apiClient.RequestWithBody("PUT", fmt.Sprintf("/v1/users/%v", authUserHashID), body)
		//
		// 		convey.Convey("Then the User should be updated", func() {
		// 			var userResp struct {
		// 				HashID       string `json:"id"`
		// 				EmailAddress string `json:"emailAddress"`
		// 				FirstName    string `json:"firstName"`
		// 				LastName     string `json:"lastName"`
		// 			}
		// 			apiClient.UnmarshalTo(&userResp)
		//
		// 			convey.So(userResp.HashID, convey.ShouldEqual, authUserHashID)
		// 			convey.So(userResp.EmailAddress, convey.ShouldEqual, "foo@bar.com")
		// 			convey.So(userResp.FirstName, convey.ShouldEqual, "Foo2")
		// 			convey.So(userResp.LastName, convey.ShouldEqual, "Bar2")
		// 		})
		// 	})
		// })
	})

}
