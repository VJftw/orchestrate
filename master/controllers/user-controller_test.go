package controllers

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/managers"
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
		userController := UserController{
			ModelManager:  modelManager,
			UserValidator: userValidator,
		}

		convey.Convey("When there is a new valid POST request", func() {
			request := http.Request{}
			writer := MockResponseWriter{}

			// modelManager.EXPECT().Save(arg0)

			userController.postHandler(writer, &request)

		})
	})
}
