package cadetGroup_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup"
)

func TestValidator(t *testing.T) {

	convey.Convey("Given a Validator", t, func() {
		validator := cadetGroup.NewValidator()

		convey.Convey("When it receives a valid name and configuration", func() {
			cG := &cadetGroup.CadetGroup{
				Name:          "Project Foo",
				Configuration: "aaaa",
			}

			convey.Convey("It should return true", func() {
				res := validator.Validate(cG)
				convey.So(res, convey.ShouldBeTrue)
			})
		})

		convey.Convey("When it receives an invalid name or configuration", func() {
			cG := &cadetGroup.CadetGroup{
				Name: "",
			}

			convey.Convey("It should return false", func() {
				res := validator.Validate(cG)
				convey.So(res, convey.ShouldBeFalse)
			})
		})
	})
}
