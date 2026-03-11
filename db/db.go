package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB //Pool of connections

func ConnectDB() {
	connStr := "user=postgres password=postgres host=localhost port=5432 dbname=builderwire_db sslmode=disable"

	db, err := sql.Open("postgres", connStr) //Opens a lazy connection to database
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping() //Tests and makes an actual connection to the database
	if err != nil {
		log.Fatal("Database unreachable: ", err)
	}

	DB = db

	log.Println("Connected to database.")

}
