package managers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/mocks"
	"github.com/vjftw/orchestrate/master/models"
)

// TestUserManager - User Manager Tests
func TestUserManager(t *testing.T) {
	persisterMock := mocks.Persister{}
	userManager := UserManager{ORM: &persisterMock}

	Convey("Given a User Model", t, func() {
		user := models.User{
			EmailAddress: "foo@bar.com",
			Password:     "axcbde3324addf3",
			FirstName:    "Foo",
			LastName:     "Bar",
		}

		Convey("When the User is Saved", func() {
			persisterMock.On("Save", &user).Return(true)
			userManager.Save(&user)
			Convey("The User should be saved on the database", func() {
				So(persisterMock.AssertCalled(t, "Save", &user), ShouldBeTrue)
			})
		})

	})

}
