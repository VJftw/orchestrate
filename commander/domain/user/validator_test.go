package user_test

import (
	"errors"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/commander/domain/user"
	"github.com/vjftw/orchestrate/commander/domain/user/mocks"
)

func TestValidator(t *testing.T) {
	convey.Convey("Given a User Validator", t, func() {
		userManager := &mocks.Manager{}
		userValidator := user.UserValidator{
			UserManager: userManager,
		}

		convey.Convey("When it recieves an invalid email address", func() {
			user := user.User{
				EmailAddress: "notanemail",
				Password:     "validPassword",
			}

			res := userValidator.Validate(&user)
			convey.Convey("It should return false", func() {
				convey.So(res, convey.ShouldBeFalse)
			})
		})

		convey.Convey("When it recieves an invalid password", func() {
			user := user.User{
				EmailAddress: "foo@bar.com",
				Password:     "iacd",
			}

			res := userValidator.Validate(&user)
			convey.Convey("It should return false", func() {
				convey.So(res, convey.ShouldBeFalse)
			})
		})

		convey.Convey("When it doesn't recieve a password", func() {
			user := user.User{
				EmailAddress: "foo@bar.com",
			}

			res := userValidator.Validate(&user)
			convey.Convey("It should return false", func() {
				convey.So(res, convey.ShouldBeFalse)
			})
		})

		convey.Convey("When it doesn't recieve an email address", func() {
			user := user.User{
				Password: "foobar1234",
			}

			res := userValidator.Validate(&user)
			convey.Convey("It should return false", func() {
				convey.So(res, convey.ShouldBeFalse)
			})
		})

		convey.Convey("When it receives a first name with non-alpha characters", func() {
			user := user.User{
				EmailAddress: "foo@bar.com",
				Password:     "foobar1234",
				FirstName:    "123123foo",
			}

			res := userValidator.Validate(&user)
			convey.Convey("It should return false", func() {
				convey.So(res, convey.ShouldBeFalse)
			})
		})

		convey.Convey("When it receives a last name with non-alpha characters", func() {
			user := user.User{
				EmailAddress: "foo@bar.com",
				Password:     "foobar1234",
				LastName:     "123123bar",
			}

			res := userValidator.Validate(&user)
			convey.Convey("It should return false", func() {
				convey.So(res, convey.ShouldBeFalse)
			})
		})

		convey.Convey("When it receives a duplicate email address", func() {
			user := user.User{
				EmailAddress: "duplicate@bar.com",
				Password:     "foobar1234",
				FirstName:    "Foo",
				LastName:     "Bar",
			}

			userManager.On("FindByEmailAddress", "duplicate@bar.com").Return(nil, nil).Once()
			res := userValidator.Validate(&user)
			convey.Convey("It should return false", func() {
				convey.So(res, convey.ShouldBeFalse)
				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})

		convey.Convey("When it receives a valid user", func() {
			user := user.User{
				EmailAddress: "foo@bar.com",
				Password:     "foobar1234",
				FirstName:    "Foo",
				LastName:     "Bar",
			}

			userManager.On("FindByEmailAddress", "foo@bar.com").Return(nil, errors.New("Not found")).Once()
			res := userValidator.Validate(&user)
			convey.Convey("It should return true", func() {
				convey.So(res, convey.ShouldBeTrue)
				convey.So(userManager.AssertExpectations(t), convey.ShouldBeTrue)
			})
		})
	})
}
