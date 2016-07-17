package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/messages"
	"github.com/vjftw/orchestrate/master/mocks"
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/routers"
)

func TestUserController(t *testing.T) {

	convey.Convey("Given a User Controller", t, func() {
		ctrl := gomock.NewController(t)
		modelManager := mocks.NewMockManager(ctrl)
		userValidator := mocks.NewMockValidator(ctrl)
		userProvider := mocks.NewMockIUserProvider(ctrl)
		userResolver := mocks.NewMockIUserResolver(ctrl)
		defer ctrl.Finish()

		userController := UserController{
			ModelManager:  modelManager,
			UserValidator: userValidator,
			UserProvider:  userProvider,
			UserResolver:  userResolver,
		}

		router := routers.NewMuxRouter([]routers.Routable{&userController}, false)

		convey.Convey("When I send a valid POST request to /v1/users", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "foobar1234",
			})

			request, _ := http.NewRequest("POST", "/v1/users", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			userValidator.EXPECT().Validate(&user).Times(1).Return(true, nil)
			modelManager.EXPECT().Save(&user).Times(1)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 201 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 201)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				var jsonResp map[string]string
				json.Unmarshal(writer.Body.Bytes(), &jsonResp)

				convey.So(jsonResp["uuid"], convey.ShouldNotBeEmpty)
			})
		})

		convey.Convey("When I send an invalid POST request to /v1/users", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "iamnotanemail",
				"password":     "a",
			})

			request, _ := http.NewRequest("POST", "/v1/users", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			vM := messages.ValidationMessage{}
			userValidator.EXPECT().Validate(&user).Times(1).Return(false, &vM)
			modelManager.EXPECT().Save(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send a malformed POST request to /v1/users", func() {
			b := []byte("malformedJSON")

			request, _ := http.NewRequest("POST", "/v1/users", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(errors.New("Malformed"))
			userValidator.EXPECT().Validate(&user).Times(0)
			modelManager.EXPECT().Save(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send an authenticated valid PUT request to /v1/users/{id}", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "bar@foo.com",
				"password":     "barfoo4321",
				"firstName":    "bar",
				"lastName":     "foo",
			})

			request, _ := http.NewRequest("PUT", "/v1/users/abcdef1234", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			modelManager.EXPECT().GetInto(&user, "uuid = ?", "abcdef1234").Times(1).Return(nil)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			user.Password = "barfoo4321"
			userValidator.EXPECT().Validate(&user).Times(1).Return(true, nil)
			modelManager.EXPECT().Save(&user).Times(1)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 200 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 200)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send an unauthenticated PUT request to /v1/users/{id}", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "bar@foo.com",
				"password":     "barfoo4321",
				"firstName":    "bar",
				"lastName":     "foo",
			})

			request, _ := http.NewRequest("PUT", "/v1/users/abcdef1234", ioutil.NopCloser(bytes.NewReader(b)))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(0)
			modelManager.EXPECT().GetInto(&user, "uuid = ?", "abcdef1234").Times(0)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(0)
			userValidator.EXPECT().Validate(&user).Times(0)
			modelManager.EXPECT().Save(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 401 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 401)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send an authenticated valid PUT request to /v1/users/{id} for a different user id", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "bar@foo.com",
				"password":     "barfoo4321",
				"firstName":    "bar",
				"lastName":     "foo",
			})

			request, _ := http.NewRequest("PUT", "/v1/users/anotherUser", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(0)
			modelManager.EXPECT().GetInto(&user, "uuid = ?", "anotherUser").Times(0)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(0)
			userValidator.EXPECT().Validate(&user).Times(0)
			modelManager.EXPECT().Save(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 403 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 403)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send a PUT request to /v1/users/{id} where the id is missing", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "bar@foo.com",
				"password":     "barfoo4321",
				"firstName":    "bar",
				"lastName":     "foo",
			})

			request, _ := http.NewRequest("PUT", "/v1/users/missing", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("missing"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			modelManager.EXPECT().GetInto(&user, "uuid = ?", "missing").Times(1).Return(errors.New("Not Found."))
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(0)
			userValidator.EXPECT().Validate(&user).Times(0)
			modelManager.EXPECT().Save(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 404 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 404)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send a malformed PUT request to /v1/users/{id}", func() {
			b := []byte("malformedJSON")

			request, _ := http.NewRequest("PUT", "/v1/users/abcdef1234", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			modelManager.EXPECT().GetInto(&user, "uuid = ?", "abcdef1234").Times(1).Return(nil)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(errors.New("Malformed JSON"))
			userValidator.EXPECT().Validate(&user).Times(0)
			modelManager.EXPECT().Save(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I send an invalid PUT request to /v1/users/{id}", func() {
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "bar",
				"password":     "b",
				"firstName":    "bar",
				"lastName":     "foo",
			})
			request, _ := http.NewRequest("PUT", "/v1/users/abcdef1234", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			modelManager.EXPECT().GetInto(&user, "uuid = ?", "abcdef1234").Times(1).Return(nil)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			vM := messages.ValidationMessage{}
			userValidator.EXPECT().Validate(&user).Times(1).Return(false, &vM)
			modelManager.EXPECT().Save(&user).Times(0)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

	})
}
