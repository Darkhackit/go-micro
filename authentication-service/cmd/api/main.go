package main

import (
	"database/sql"
	"github.com/Darkhackit/go-micro-authentication/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication server...")

	conn := connectToDB()

	if conn == nil {
		log.Fatal("Error connecting to database")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
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

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready")
			count++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}
		if count > 10 {
			log.Println("Too many retries")
			return nil
		}
		log.Println("Retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
