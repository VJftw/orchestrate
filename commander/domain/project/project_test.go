package project_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/project"
)

func TestProject(t *testing.T) {
	convey.Convey("Given a Project", t, func() {
		p := project.Project{}

		convey.Convey("It should return the UUID", func() {
			convey.So(p.GetUUID(), convey.ShouldBeEmpty)
		})
	})
}
