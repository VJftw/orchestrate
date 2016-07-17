package persisters

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/vjftw/orchestrate/master/models"
)

// GORMPersister - Persistence using the GORM library
type GORMPersister struct {
	db *gorm.DB
}

func NewGORMPersister() *GORMPersister {
	gormPersister := GORMPersister{}

	db, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

	gormPersister.db = db

	return &gormPersister
}

// Save - saves an object using GORM
func (gP GORMPersister) Save(v models.Model) error {
	gP.db.Save(v)
	return nil
}

func (gP GORMPersister) GetInto(
	v models.Model,
	query interface{},
	args ...interface{},
) error {
	gP.db.Where(query, args).First(v)
	return nil
}

func (gP GORMPersister) Delete(m models.Model) error {
	gP.db.Delete(m, nil)
	return nil
}
