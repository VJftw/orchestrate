package ephemeral

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAuthToken(t *testing.T) {
	convey.Convey("Given an ephemeral auth token", t, func() {
		authToken := AuthToken{
			Token: "abcdef1234",
		}

		convey.Convey("It should return a serializable map", func() {
			convey.So(authToken.ToMap()["authToken"], convey.ShouldEqual, "abcdef1234")
		})
	})
}
