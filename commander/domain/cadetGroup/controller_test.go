package cadetGroup_test

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
	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup/mocks"
	"github.com/vjftw/orchestrate/commander/domain/project"
	project_mocks "github.com/vjftw/orchestrate/commander/domain/project/mocks"
	"github.com/vjftw/orchestrate/commander/domain/user"
	user_mocks "github.com/vjftw/orchestrate/commander/domain/user/mocks"
	"github.com/vjftw/orchestrate/commander/routers"
)

func TestController(t *testing.T) {
	convey.Convey("Given a Controller", t, func() {
		userManager := &user_mocks.Manager{}
		projectManager := &project_mocks.Manager{}
		cadetGroupResolver := &mocks.Resolver{}
		cadetGroupValidator := &mocks.Validator{}
		cadetGroupManager := &mocks.Manager{}

		cadetGroupController := cadetGroup.Controller{
			UserManager:         userManager,
			ProjectManager:      projectManager,
			CadetGroupResolver:  cadetGroupResolver,
			CadetGroupValidator: cadetGroupValidator,
			CadetGroupManager:   cadetGroupManager,
		}

		router := routers.NewMuxRouter([]routers.Routable{&cadetGroupController}, false)

		convey.Convey("When I send an authenticated valid POST request to /v1/projects/{project-id}/cadetGroups", func() {

			b, _ := json.Marshal(map[string]string{
				"name": "Foo Pre-emptible workers",
			})

			request, _ := http.NewRequest("POST", "/v1/projects/a-project/cadetGroups", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("testUser"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "testUser").Return(user, nil).Once()
			project := &project.Project{}
			projectManager.On("FindByUserAndUUID", user, "a-project").Return(project, nil).Once()
			cG := &cadetGroup.CadetGroup{}
			cadetGroupManager.On("NewForProject", project).Return(cG).Once()
			cadetGroupResolver.On("FromRequest", cG, request.Body).Return(nil).Once()
			cadetGroupValidator.On("Validate", cG).Return(true).Once()
			cadetGroupManager.On("Save", cG).Return(nil).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 201 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 201)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				var jsonResp map[string]interface{}
				json.Unmarshal(writer.Body.Bytes(), &jsonResp)

				convey.So(jsonResp["uuid"], convey.ShouldNotBeEmpty)
				convey.So(jsonResp["key"], convey.ShouldNotBeEmpty)
				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send an authenticated invalid POST request to /v1/projects/{project-id}/cadetGroups", func() {

			b, _ := json.Marshal(map[string]string{
				"name": "a",
			})

			request, _ := http.NewRequest("POST", "/v1/projects/a-project/cadetGroups", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("testUser"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "testUser").Return(user, nil).Once()
			project := &project.Project{}
			projectManager.On("FindByUserAndUUID", user, "a-project").Return(project, nil).Once()
			cG := &cadetGroup.CadetGroup{}
			cadetGroupManager.On("NewForProject", project).Return(cG).Once()
			cadetGroupResolver.On("FromRequest", cG, request.Body).Return(nil).Once()
			cadetGroupValidator.On("Validate", cG).Return(false).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send an authenticated malformed POST request to /v1/projects/{project-id}/cadetGroups", func() {

			b := []byte("malformedJSON")

			request, _ := http.NewRequest("POST", "/v1/projects/a-project/cadetGroups", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("testUser"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "testUser").Return(user, nil).Once()
			project := &project.Project{}
			projectManager.On("FindByUserAndUUID", user, "a-project").Return(project, nil).Once()
			cG := &cadetGroup.CadetGroup{}
			cadetGroupManager.On("NewForProject", project).Return(cG).Once()
			cadetGroupResolver.On("FromRequest", cG, request.Body).Return(errors.New("Malformed JSON")).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send a falsely authenticated POST request to /v1/projects/{project-id}/cadetGroups", func() {

			b, _ := json.Marshal(map[string]string{
				"name": "Foo Pre-emptible workers",
			})

			request, _ := http.NewRequest("POST", "/v1/projects/a-project/cadetGroups", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("unknownUser"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			userManager.On("FindByUUID", "unknownUser").Return(nil, errors.New("Not found")).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 401 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 401)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send an authenticated valid POST request to /v1/projects/{project-id}/cadetGroups where the project-id is missing", func() {

			b, _ := json.Marshal(map[string]string{
				"name": "Foo Pre-emptible workers",
			})

			request, _ := http.NewRequest("POST", "/v1/projects/a-project/cadetGroups", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("testUser"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "testUser").Return(user, nil).Once()
			projectManager.On("FindByUserAndUUID", user, "a-project").Return(nil, errors.New("Not found")).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 403 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 403)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send an authenticated GET request to /v1/projects/{project-id}/cadetGroups", func() {

			request, _ := http.NewRequest("GET", "/v1/projects/a-project/cadetGroups", nil)
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("testUser"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "testUser").Return(user, nil).Once()
			cadetGroups := &[]cadetGroup.CadetGroup{}
			cadetGroupManager.On("FindByUserAndProjectUUID", user, "a-project").Return(cadetGroups, nil).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 200 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 200)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				var jsonResp map[string]interface{}
				json.Unmarshal(writer.Body.Bytes(), &jsonResp)

				convey.So(len(jsonResp), convey.ShouldEqual, 0)

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send a falsely authenticated GET request to /v1/projects/{project-id}/cadetGroups", func() {

			request, _ := http.NewRequest("GET", "/v1/projects/a-project/cadetGroups", nil)
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("unknownUser"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			userManager.On("FindByUUID", "unknownUser").Return(nil, errors.New("Not found")).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 401 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 401)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(cadetGroupValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

	})
}
