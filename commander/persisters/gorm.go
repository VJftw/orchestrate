package persisters

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// SQLite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// GORM - Persistence using the GORM library
type GORM struct {
	db *gorm.DB
}

// NewGORM - Initialises a connection to a GORM storage
func NewGORM(models ...interface{}) *GORM {
	gormPersister := GORM{}

	db, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(models...)

	gormPersister.db = db

	return &gormPersister
}

// Save - saves an object using GORM
func (gP GORM) Save(v Persistable) error {
	gP.db.Save(v)
	return nil
}

// GetInto - Searches the Storage and places the result into a given Model based on the given query
func (gP GORM) GetInto(
	v Persistable,
	query interface{},
	args ...interface{},
) error {
	gP.db.Where(query, args...).First(v)
	return nil
}

// Delete - Deletes a given model from the storage
func (gP GORM) Delete(m Persistable) error {
	gP.db.Delete(m, nil)
	return nil
}
