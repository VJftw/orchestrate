package cadetGroup_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup"
)

func TestCadetGroup(t *testing.T) {
	convey.Convey("Given a Cadet Group", t, func() {
		cG := cadetGroup.CadetGroup{}

		convey.Convey("It should return the UUID", func() {
			convey.So(cG.GetUUID(), convey.ShouldBeEmpty)
		})
	})
}
