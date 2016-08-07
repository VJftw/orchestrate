package project_test

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

	"github.com/dgrijalva/jwt-go"
	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/project"
	"github.com/vjftw/orchestrate/commander/domain/project/mocks"
	"github.com/vjftw/orchestrate/commander/domain/user"
	user_mocks "github.com/vjftw/orchestrate/commander/domain/user/mocks"
	"github.com/vjftw/orchestrate/commander/routers"
)

func TestController(t *testing.T) {
	convey.Convey("Given a Project Controller", t, func() {
		userManager := &user_mocks.Manager{}
		projectManager := &mocks.Manager{}
		projectResolver := &mocks.Resolver{}
		projectValidator := &mocks.Validator{}

		projectController := project.Controller{
			UserManager:      userManager,
			ProjectManager:   projectManager,
			ProjectResolver:  projectResolver,
			ProjectValidator: projectValidator,
		}

		router := routers.NewMuxRouter([]routers.Routable{&projectController}, false)

		convey.Convey("When I send an authenticated valid POST request to /v1/projects", func() {
			b, _ := json.Marshal(map[string]string{
				"name": "Project Foo",
			})

			request, _ := http.NewRequest("POST", "/v1/projects", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "abcdef1234").Return(user, nil).Once()
			p := &project.Project{}
			projectManager.On("NewForUser", user).Return(p).Once()
			projectResolver.On("FromRequest", p, request.Body).Return(nil).Once()
			projectValidator.On("Validate", p).Return(true).Once()
			projectManager.On("Save", p).Return(nil).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 201 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 201)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				var jsonResp map[string]interface{}
				json.Unmarshal(writer.Body.Bytes(), &jsonResp)

				convey.So(jsonResp["uuid"], convey.ShouldNotBeEmpty)
				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send an authenticated invalid POST request to /v1/projects", func() {
			b, _ := json.Marshal(map[string]string{
				"name": "a",
			})

			request, _ := http.NewRequest("POST", "/v1/projects", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "abcdef1234").Return(user, nil).Once()
			p := &project.Project{}
			projectManager.On("NewForUser", user).Return(p).Once()
			projectResolver.On("FromRequest", p, request.Body).Return(nil).Once()
			projectValidator.On("Validate", p).Return(false).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send an authenticated malformed POST request to /v1/projects", func() {
			b := []byte("malformedJSON")

			request, _ := http.NewRequest("POST", "/v1/projects", ioutil.NopCloser(bytes.NewReader(b)))
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "abcdef1234").Return(user, nil).Once()
			p := &project.Project{}
			projectManager.On("NewForUser", user).Return(p).Once()
			projectResolver.On("FromRequest", p, request.Body).Return(errors.New("Malformed JSON")).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 400 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 400)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectResolver.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectValidator.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When I send a falsely authenticated POST request to /v1/projects", func() {
			b, _ := json.Marshal(map[string]string{
				"name": "Project Foo",
			})

			request, _ := http.NewRequest("POST", "/v1/projects", ioutil.NopCloser(bytes.NewReader(b)))
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
			})
		})

		convey.Convey("When I send an authenticated GET request to /v1/projects", func() {
			request, _ := http.NewRequest("GET", "/v1/projects", nil)
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "abcdef1234").Return(user, nil).Once()
			p := &[]project.Project{}
			projectManager.On("FindByUser", user).Return(p).Once()

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 200 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 200)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				var jsonResp map[string]interface{}
				json.Unmarshal(writer.Body.Bytes(), &jsonResp)

				convey.So(len(jsonResp), convey.ShouldEqual, 0)
			})
		})

		convey.Convey("When I send a falsely authenticated GET request to /v1/projects", func() {
			request, _ := http.NewRequest("GET", "/v1/projects", nil)
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
			})
		})

		convey.Convey("When I send an authenticated GET request to /v1/projects/{projectUUID} where I own the project", func() {
			request, _ := http.NewRequest("GET", "/v1/projects/abcd-ef-213-234a", nil)
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "abcdef1234").Return(user, nil).Once()
			project := &project.Project{
				Name: "Project Foo",
			}
			projectManager.On("FindByUserAndUUID", user, "abcd-ef-213-234a").Return(project, nil)

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 200 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 200)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)

				var jsonResp map[string]interface{}
				json.Unmarshal(writer.Body.Bytes(), &jsonResp)

				convey.So(jsonResp["name"], convey.ShouldEqual, "Project Foo")
			})
		})

		convey.Convey("When I send a falsely authenticated GET request to /v1/projects/{projectUUID}", func() {
			request, _ := http.NewRequest("GET", "/v1/projects/abcd-ef-213-234a", nil)
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
			})
		})

		convey.Convey("When I send an authenticated GET request to /v1/projects/{projectUUID} where I do not own the project, or it does not exist", func() {
			request, _ := http.NewRequest("GET", "/v1/projects/unknownProject", nil)
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userUUID": string("abcdef1234"),
				"nbf":      time.Now().Unix(),
			})
			jwtTokenStr, _ := jwtToken.SignedString([]byte("hmacSecret"))
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtTokenStr))
			writer := httptest.NewRecorder()

			user := &user.User{}
			userManager.On("FindByUUID", "abcdef1234").Return(user, nil).Once()
			projectManager.On("FindByUserAndUUID", user, "unknownProject").Return(nil, errors.New("Not found"))

			router.Handler.ServeHTTP(writer, request)

			convey.Convey("Then it should give the correct 404 response", func() {
				convey.So(writer.Code, convey.ShouldEqual, 404)
				convey.So(writer.Header().Get("Content-type"), convey.ShouldEqual, "application/json; charset=UTF-8")

				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
				convey.So(projectManager.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

	})
}
