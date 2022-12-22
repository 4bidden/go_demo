package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "messages"
)

var dbPool = sync.Pool{
	New: func() interface{} {
		// Connect to the database without sslmode

		// Create a connection string
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := gorm.Open("postgres", psqlInfo)
		if err != nil {
			// Handle error
			log.Fatal("Error connecting to database: ", err)
		}
		db.LogMode(false)
		return db
	},
}

// GetConnection returns a connection from the pool
func GetConnection() *gorm.DB {
	return dbPool.Get().(*gorm.DB)
}

// PutConnection returns a connection to the pool
func PutConnection(db *gorm.DB) {
	dbPool.Put(db)
}

// create a function to migrate the database and message table
func Migrate() {
	log.Println("Migrating database")

	db := GetConnection()
	defer PutConnection(db)
	db.AutoMigrate(&Message{})
}
