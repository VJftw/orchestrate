package cadetGroup_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup"
)

func TestResolver(t *testing.T) {
	convey.Convey("Given a CadetGroup Resolver", t, func() {
		cadetGroupResolver := cadetGroup.NewResolver()

		convey.Convey("When it receives a name and configuration", func() {
			b, _ := json.Marshal(map[string]string{
				"name":          "Foo Pre-emptible workers",
				"configuration": "blahblah",
			})
			body := ioutil.NopCloser(bytes.NewReader(b))

			cG := &cadetGroup.CadetGroup{}

			err := cadetGroupResolver.FromRequest(cG, body)

			convey.Convey("It should resolve correctly", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(cG.Name, convey.ShouldEqual, "Foo Pre-emptible workers")
				convey.So(cG.Configuration, convey.ShouldEqual, "blahblah")
			})
		})

		convey.Convey("When it doesn't receive a name", func() {
			b, _ := json.Marshal(map[string]string{
				"configuration": "blahblah",
			})
			body := ioutil.NopCloser(bytes.NewReader(b))

			cG := &cadetGroup.CadetGroup{}

			err := cadetGroupResolver.FromRequest(cG, body)

			convey.Convey("It should return an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(cG.Name, convey.ShouldBeEmpty)
				convey.So(cG.Configuration, convey.ShouldBeEmpty)
			})
		})

		convey.Convey("When it doesn't receive a configuration", func() {
			b, _ := json.Marshal(map[string]string{
				"name": "blahblah",
			})
			body := ioutil.NopCloser(bytes.NewReader(b))

			cG := &cadetGroup.CadetGroup{}

			err := cadetGroupResolver.FromRequest(cG, body)

			convey.Convey("It should return an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(cG.Name, convey.ShouldBeEmpty)
				convey.So(cG.Configuration, convey.ShouldBeEmpty)
			})
		})

		convey.Convey("When it receives malformed JSON", func() {
			b := []byte("malformedJSON")
			body := ioutil.NopCloser(bytes.NewReader(b))
			cG := &cadetGroup.CadetGroup{}
			err := cadetGroupResolver.FromRequest(cG, body)

			convey.Convey("It should return an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(cG.Name, convey.ShouldBeEmpty)
				convey.So(cG.Configuration, convey.ShouldBeEmpty)
			})
		})
	})
}
