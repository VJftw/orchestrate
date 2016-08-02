package project_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/project"
)

func TestResolver(t *testing.T) {
	convey.Convey("Given a Project Resolver", t, func() {
		projectResolver := project.ProjectResolver{}

		convey.Convey("When it receives a name", func() {
			b, _ := json.Marshal(map[string]string{
				"name": "Project Foo",
			})
			body := ioutil.NopCloser(bytes.NewReader(b))

			p := &project.Project{}

			err := projectResolver.FromRequest(p, body)

			convey.Convey("It should resolve correctly", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(p.Name, convey.ShouldEqual, "Project Foo")
			})

		})
		convey.Convey("When it doesn't receive a name", func() {
			b, _ := json.Marshal(map[string]string{
				"foo": "bar",
			})
			body := ioutil.NopCloser(bytes.NewReader(b))

			p := &project.Project{}

			err := projectResolver.FromRequest(p, body)

			convey.Convey("It should return an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(p.Name, convey.ShouldBeEmpty)
			})
		})
		convey.Convey("When it receives malformed JSON", func() {
			b := []byte("malformedJSON")
			body := ioutil.NopCloser(bytes.NewReader(b))
			p := &project.Project{}
			err := projectResolver.FromRequest(p, body)

			convey.Convey("It should return an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(p.Name, convey.ShouldBeEmpty)
			})
		})
	})
}
