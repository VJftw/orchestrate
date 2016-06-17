package controllers

import (
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

// TestUserController - User Controller Tests
func TestUserController(t *testing.T) {

	userController := UserController{}

	Convey("Given a Router", t, func() {
		router := mux.Router{}

		userController.AddRoutes(&router)

		Convey("The route: POST /v1/users should be added", func() {
			path, _ := router.Get("POSTusers").GetPathTemplate()

			So(path, ShouldEqual, "/v1/users")
		})
	})
}
