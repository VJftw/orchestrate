package user_test

import (
	"errors"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/user"
	"github.com/vjftw/orchestrate/commander/persisters/mocks"
)

func TestManager(t *testing.T) {
	convey.Convey("Given a Manager", t, func() {
		gormPersister := &mocks.Persister{}
		userManager := user.NewManager(gormPersister)

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

		convey.Convey("It should get a User into an object if it exists", func() {
			user := user.User{}

			gormPersister.On("GetInto", &user, "email_address = ?", []interface{}{"foo@bar.com"}).Return(nil).Once()

			convey.So(userManager.GetInto(&user, "email_address = ?", "foo@bar.com"), convey.ShouldBeNil)
			convey.So(gormPersister.AssertExpectations(t), convey.ShouldBeTrue)
		})

		convey.Convey("It should return errors if necessary", func() {
			user := user.User{}

			gormPersister.On("GetInto", &user, "email_address = ?", []interface{}{"foo@bar.com"}).Return(errors.New("invalid column")).Once()

			err := userManager.GetInto(&user, "email_address = ?", "foo@bar.com")
			convey.So(err, convey.ShouldNotBeNil)

			convey.So(gormPersister.AssertExpectations(t), convey.ShouldBeTrue)
		})

	})
}
