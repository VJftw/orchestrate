package persisters

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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

	gP.DB = db
}

// Save - saves an object using GORM
func (gP *GORMPersister) Save(v interface{}) {
	gP.DB.Save(v)
}
