package user_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/user"
	"github.com/vjftw/orchestrate/commander/persisters/mocks"
)

func TestManager(t *testing.T) {
	convey.Convey("Given a Manager", t, func() {
		gormPersister := &mocks.Persister{}
		userManager := user.UserManager{
			GORMPersister: gormPersister,
		}

		convey.Convey("It should Save a given User", func() {
			user := user.User{}
			gormPersister.On("Save", &user).Return(nil).Once()

			convey.So(userManager.Save(&user), convey.ShouldBeNil)

			convey.So(gormPersister.AssertExpectations(t), convey.ShouldBeTrue)
		})

		convey.Convey("It should Delete a given User", func() {
			user := user.User{}
			gormPersister.On("Delete", &user).Return(nil).Once()

			convey.So(userManager.Delete(&user), convey.ShouldBeNil)

			convey.So(gormPersister.AssertExpectations(t), convey.ShouldBeTrue)
		})

	})
}
