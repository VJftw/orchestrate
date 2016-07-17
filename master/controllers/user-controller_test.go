package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/unrolled/render"
	"github.com/vjftw/orchestrate/master/messages"
	"github.com/vjftw/orchestrate/master/mocks"
	"github.com/vjftw/orchestrate/master/models"
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
			render:        render.New(),
			ModelManager:  modelManager,
			UserValidator: userValidator,
			UserProvider:  userProvider,
			UserResolver:  userResolver,
		}

		convey.Convey("When there is a valid POST request", func() {
			request := http.Request{}
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "foo@bar.com",
				"password":     "foobar1234",
			})
			request.Body = ioutil.NopCloser(bytes.NewReader(b))
			rw := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			modelManager.EXPECT().Save(&user).Times(1)
			userValidator.EXPECT().Validate(&user).Times(1).Return(true, nil)

			userController.postHandler(rw, &request)

			convey.Convey("Then it should give the correct 201 response", func() {
				convey.So(rw.Code, convey.ShouldEqual, 201)
				convey.So(rw.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				var jsonResp map[string]string
				json.Unmarshal(rw.Body.Bytes(), &jsonResp)

				convey.So(jsonResp["uuid"], convey.ShouldNotBeEmpty)
			})
		})

		convey.Convey("When there is an invalid POST request", func() {
			request := http.Request{}
			b, _ := json.Marshal(map[string]string{
				"emailAddress": "iamnotanemail",
				"password":     "a",
			})
			request.Body = ioutil.NopCloser(bytes.NewReader(b))
			rw := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(nil)
			vM := messages.ValidationMessage{}
			userValidator.EXPECT().Validate(&user).Times(1).Return(false, &vM)
			modelManager.EXPECT().Save(&user).Times(0)

			userController.postHandler(rw, &request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(rw.Code, convey.ShouldEqual, 400)
				convey.So(rw.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When there is a malformed POST request", func() {
			request := http.Request{}
			b := []byte("aasdasds")
			request.Body = ioutil.NopCloser(bytes.NewReader(b))
			rw := httptest.NewRecorder()

			user := models.User{}
			userProvider.EXPECT().New().Times(1).Return(&user)
			userResolver.EXPECT().FromRequest(&user, request.Body).Times(1).Return(errors.New("Malformed"))
			userValidator.EXPECT().Validate(&user).Times(0)
			modelManager.EXPECT().Save(&user).Times(0)

			userController.postHandler(rw, &request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(rw.Code, convey.ShouldEqual, 400)
				convey.So(rw.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

	})
}
