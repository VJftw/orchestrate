package persisters

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	// postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewGORMDB - Initialises a connection to a GORM storage
func NewGORMDB(models ...interface{}) *gorm.DB {

	if !waitForService(fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), "5432")) {
		panic("Could not find database")
	}

	// db, err := gorm.Open("sqlite3", "test.db")
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	))

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(models...)

	return db
}
