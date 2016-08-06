package project_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/project"
)

func TestValidator(t *testing.T) {

	convey.Convey("Given a Validator", t, func() {
		validator := project.ProjectValidator{}

		convey.Convey("When it receives a valid name", func() {
			p := &project.Project{
				Name: "Project Foo",
			}

			convey.Convey("It should return true", func() {
				res := validator.Validate(p)
				convey.So(res, convey.ShouldBeTrue)
			})
		})

		convey.Convey("When it receives an invalid name", func() {
			p := &project.Project{
				Name: "",
			}

			convey.Convey("It should return false", func() {
				res := validator.Validate(p)
				convey.So(res, convey.ShouldBeFalse)
			})
		})
	})
}
