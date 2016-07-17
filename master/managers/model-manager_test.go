package managers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/vjftw/orchestrate/master/mocks"
	"github.com/vjftw/orchestrate/master/models"
)

func TestModelManager(t *testing.T) {
	convey.Convey("Given a Model Manager", t, func() {
		ctrl := gomock.NewController(t)
		persister := mocks.NewMockPersister(ctrl)

		modelManager := ModelManager{
			GORMPersister: persister,
		}

		convey.Convey("Then I can call Save() with a Model", func() {
			model := models.User{}

			persister.EXPECT().Save(&model).Times(1)
			modelManager.Save(&model)
		})

		convey.Convey("Then I can call Delete() with a Model", func() {
			model := models.User{}

			persister.EXPECT().Delete(&model).Times(1)
			modelManager.Delete(&model)
		})
	})
}
