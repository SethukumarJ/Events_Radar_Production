package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectDB() *sql.DB {

	//loading parameters from env file
	// username := os.Getenv("DB_NAME")
	// password := os.Getenv("DB_PASSWORD")
	// databaseHost := os.Getenv("DB_HOST")
	databaseName := os.Getenv("DB_NAME")
	//formatting
	// dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)

	dbURI := os.Getenv("DB_SOURCE")

	//open connection to database

	db, err := sql.Open("postgres", dbURI)

	if err != nil {
		log.Fatal(err)
	}

	// verify connectioin to the database is still alive

	err = db.Ping()
	if err != nil {
		fmt.Println("error in pinging database")
		log.Fatal(err)
	}

	log.Println("/nConnected to database", databaseName)

	return db

}
