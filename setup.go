package main

import (
	"os"
	"fmt"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/calvintran1478/api-doc-server/utils"
)

func main() {
	// Load environment variables
	utils.LoadEnv(".env")

	// Set up database connection
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_CONN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Create user table
	userTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			user_id UUID PRIMARY KEY,
			email VARCHAR UNIQUE,
			password VARCHAR NOT NULL,
			first_name VARCHAR NOT NULL,
			last_name VARCHAR NOT NULL
		);
	`
	_, err = pool.Exec(context.Background(), userTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user table", err)
		os.Exit(1)
	}

	// Create project table
	projectTableQuery := `
		CREATE TABLE IF NOT EXISTS projects (
			project_id VARCHAR PRIMARY KEY,
			name VARCHAR UNIQUE,
			user_id UUID,
			CONSTRAINT project_user_id_fkey FOREIGN KEY(user_id) REFERENCES users(user_id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			UNIQUE (user_id, name)
		);
	`
	_, err = pool.Exec(context.Background(), projectTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating project table", err)
		os.Exit(1)
	}
}
