package main

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func main() {

	initDB()

}

func initDB() *sql.DB {

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to DB")
	}
}

func connectToDB() *sql.DB {
	counts := 0
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not ready yet...")

		} else {
			log.Println("connected to database")
			return connection
		}

		if counts > 10 {
			return nil
		}
		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
a
