package persisters

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/vjftw/orchestrate/master/models"
)

// GORMPersister - Persistence using the GORM library
type GORMPersister struct {
	DB *gorm.DB
}

// Init - Initialises a new database
func (gP *GORMPersister) Init() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

	gP.DB = db
}

// Save - saves an object using GORM
func (gP *GORMPersister) Save(v interface{}) {
	gP.DB.Save(v)
}

func (gP *GORMPersister) FindInto(
	v interface{},
	query interface{},
	args ...interface{},
) {
	gP.DB.Where(query, args).First(v)
}
