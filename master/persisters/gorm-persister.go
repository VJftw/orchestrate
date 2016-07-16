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
func (gP GORMPersister) Save(v models.Model) {
	gP.db.Save(v)
}

func (gP GORMPersister) FindInto(
	v models.Model,
	query interface{},
	args ...interface{},
) {
	gP.db.Where(query, args).First(v)
}

func (gP GORMPersister) Exists(
	e models.Model,
	query interface{},
	args ...interface{},
) bool {
	gP.FindInto(e, query, args)

	if len(e.GetUUID()) > 0 {
		return true
	}
	return false
}

func (gP GORMPersister) Delete(m models.Model) bool {
	gP.db.Delete(m, nil)
	return true
}
