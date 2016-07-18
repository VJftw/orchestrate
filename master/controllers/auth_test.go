package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/mock/managers"
	"github.com/vjftw/orchestrate/master/mock/providers"
	"github.com/vjftw/orchestrate/master/mock/resolvers"
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/models/ephemeral"
	"github.com/vjftw/orchestrate/master/routers"
)

func TestAuth(t *testing.T) {

	convey.Convey("Given an Auth Controller", t, func() {
		ctrl := gomock.NewController(t)
		userProvider := providers.NewMockIUser(ctrl)
		userResolver := resolvers.NewMockIUser(ctrl)
		modelManager := managers.NewMockIManager(ctrl)
		authTokenProvider := providers.NewMockIAuthToken(ctrl)
		defer ctrl.Finish()

		authController := Auth{
			UserProvider:      userProvider,
			UserResolver:      userResolver,
			ModelManager:      modelManager,
			AuthTokenProvider: authTokenProvider,
		}
		router := routers.NewMuxRouter([]routers.Routable{&authController}, false)

		convey.Convey("When I send valid credentials in a POST request to /v1/auth", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "foobar1234",
			})

			request, _ := http.NewRequest("POST", "/v1/auth", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			user.EmailAddress = "foo@bar.com"
			user.Password = "foobar1234"
			modelManager.EXPECT().GetInto(&user, "email_address = ?", "foo@bar.com").Times(1).Return(nil)
			user.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(user.Password), 10)
			authToken := ephemeral.AuthToken{
				Token: "newJWTToken",
			}
			authTokenProvider.EXPECT().NewFromUser(&user).Times(1).Return(&authToken, nil)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 201 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 201)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				var jsonResp map[string]string
				json.Unmarshal(writer.Body.Bytes(), &jsonResp)

				convey.So(jsonResp["authToken"], convey.ShouldEqual, "newJWTToken")
			})
		})

		convey.Convey("When I send invalid credentials in a POST request to /v1/auth", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "invalid",
			})

			request, _ := http.NewRequest("POST", "/v1/auth", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			user.EmailAddress = "foo@bar.com"
			modelManager.EXPECT().GetInto(&user, "email_address = ?", "foo@bar.com").Times(1).Return(nil)

			authTokenProvider.EXPECT().NewFromUser(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 401 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 401)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send a malformed POST request to /v1/auth", func() {
			b := []byte("malformedJSON")

			request, _ := http.NewRequest("POST", "/v1/auth", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(errors.New("Malformed"))
			modelManager.EXPECT().GetInto(&user, "email_address = ?", "foo@bar.com").Times(0)
			authTokenProvider.EXPECT().NewFromUser(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send credentials in a POST request to /v1/auth for a missing user", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "missing@missing.com",
				"password":     "barfoo4321",
			})

			request, _ := http.NewRequest("POST", "/v1/auth", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			user.EmailAddress = "missing@missing.com"
			modelManager.EXPECT().GetInto(&user, "email_address = ?", "missing@missing.com").Times(1).Return(errors.New("User not found."))
			authTokenProvider.EXPECT().NewFromUser(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 404 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 404)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When there is an error with the JWT configuration", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "barfoo4321",
			})

			request, _ := http.NewRequest("POST", "/v1/auth", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			user.EmailAddress = "foo@bar.com"
			user.Password = "foobar1234"
			modelManager.EXPECT().GetInto(&user, "email_address = ?", "foo@bar.com").Times(1).Return(nil)
			user.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(user.Password), 10)
			authTokenProvider.EXPECT().NewFromUser(&user).Times(1).Return(nil, errors.New("JWT error"))

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})
	})
}
