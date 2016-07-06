package ephemeral

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAuthEphmeral(t *testing.T) {
	convey.Convey("Given an ephemeral auth token", t, func() {
		authToken := AuthEphemeral{
			Token: "abcdef1234",
		}

		convey.Convey("It should return a serializable map", func() {
			convey.So(authToken.ToMap()["authToken"], convey.ShouldEqual, "abcdef1234")
		})
	})
}
