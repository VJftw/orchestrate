package managers

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/mock"
	"github.com/vjftw/orchestrate/master/models"
)

func TestModelManager(t *testing.T) {
	convey.Convey("Given a Model Manager", t, func() {
		ctrl := gomock.NewController(t)
		persister := mock.NewMockPersister(ctrl)

		modelManager := ModelManager{
			GORMPersister: persister,
		}

		convey.Convey("When I call Save() with a Model", func() {
			model := models.User{}

			persister.EXPECT().Save(&model).Times(1).Return(nil)
			res := modelManager.Save(&model)
			convey.Convey("Then it should save successfully", func() {
				convey.So(res, convey.ShouldBeNil)
			})
		})

		convey.Convey("When I can call Delete() with a Model", func() {
			model := models.User{}

			persister.EXPECT().Delete(&model).Times(1).Return(nil)
			res := modelManager.Delete(&model)

			convey.Convey("Then it should delete successfully", func() {
				convey.So(res, convey.ShouldBeNil)
			})
		})

		convey.Convey("When I can call GetInto() to query", func() {
			model := models.User{}
			query := "emailAddress = ?"
			args := "foo@bar.com"

			convey.Convey("Then it gets successfully", func() {
				persister.EXPECT().GetInto(&model, query, []interface{}{args}).Times(1).Return(nil)
				res := modelManager.GetInto(&model, query, args)
				convey.So(res, convey.ShouldBeNil)
			})

			convey.Convey("Or it gets unsuccessfully", func() {
				persister.EXPECT().GetInto(&model, query, []interface{}{args}).Times(1).Return(errors.New("Not Found"))
				res := modelManager.GetInto(&model, query, args)
				convey.So(res, convey.ShouldNotBeNil)
			})

		})

	})
}
