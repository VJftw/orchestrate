package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/unrolled/render"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/messages"
	"github.com/vjftw/orchestrate/master/models"
	"github.com/vjftw/orchestrate/master/providers"
	"github.com/vjftw/orchestrate/master/validators"
)

type MockResponseWriter struct {
}

func (mRW MockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (mRW MockResponseWriter) Write(d []byte) (int, error) {
	return 0, nil
}

func (mRW MockResponseWriter) WriteHeader(s int) {
}

func TestUserController(t *testing.T) {
	convey.Convey("Given a User Controller", t, func() {
		ctrl := gomock.NewController(t)
		modelManager := managers.NewMockManager(ctrl)
		userValidator := validators.NewMockValidator(ctrl)
		userProvider := providers.NewMockIUserProvider(ctrl)
		defer ctrl.Finish()
		userController := UserController{
			render:        render.New(),
			ModelManager:  modelManager,
			UserValidator: userValidator,
			UserProvider:  userProvider,
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
			modelManager.EXPECT().Save(&user).Times(1)
			userValidator.EXPECT().Validate(&user).Times(1).Return(true, nil)

			userController.postHandler(rw, &request)

			convey.Convey("Then it should give the correct 201 response", func() {
				convey.So(rw.Code, convey.ShouldEqual, 201)
				convey.So(rw.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				var jsonResp map[string]string
				json.Unmarshal(rw.Body.Bytes(), &jsonResp)

				convey.So(jsonResp["uuid"], convey.ShouldNotBeEmpty)
				convey.So(jsonResp["emailAddress"], convey.ShouldEqual, "foo@bar.com")
				convey.So(jsonResp["firstName"], convey.ShouldEqual, "")
				convey.So(jsonResp["lastName"], convey.ShouldEqual, "")
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
