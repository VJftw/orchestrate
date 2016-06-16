package persisters

import "github.com/jinzhu/gorm"

// GORMPersister - Persistence using the GORM library
type GORMPersister struct {
	DB *gorm.DB
}

func (gP *GORMPersister) init() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	gP.DB = db
}
