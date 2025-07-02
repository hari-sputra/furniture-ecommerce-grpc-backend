package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB(ctx context.Context) *sql.DB {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	DB, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	err = DB.PingContext(ctx)

	if err != nil {
		log.Fatal("Failed to ping DB", err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return DB
}
