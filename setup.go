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

	// Create endpoint table
	endpointTableQuery := `
		CREATE TABLE IF NOT EXISTS endpoints (
			endpoint_id VARCHAR PRIMARY KEY,
			method VARCHAR NOT NULL,
			path VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			project_id VARCHAR,
			CONSTRAINT endpoint_project_id_fkey FOREIGN KEY(project_id) REFERENCES projects(project_id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
		);
	`
	_, err = pool.Exec(context.Background(), endpointTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating endpoint table", err)
		os.Exit(1)
	}

	// Create path parameter table
	pathParameterTableQuery := `
		CREATE TABLE IF NOT EXISTS path_parameters (
			name VARCHAR NOT NULL,
			type VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			endpoint_id VARCHAR,
			CONSTRAINT path_parameter_endpoint_id_fkey FOREIGN KEY(endpoint_id) REFERENCES endpoints(endpoint_id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			PRIMARY KEY(name, endpoint_id)
		);
	`
	_, err = pool.Exec(context.Background(), pathParameterTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating path parameter table", err)
		os.Exit(1)
	}

	// Create query parameter table
	queryParameterTableQuery := `
		CREATE TABLE IF NOT EXISTS query_parameters (
			name VARCHAR NOT NULL,
			type VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			required BOOLEAN NOT NULL,
			endpoint_id VARCHAR,
			CONSTRAINT query_parameter_endpoint_id_fkey FOREIGN KEY(endpoint_id) REFERENCES endpoints(endpoint_id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			PRIMARY KEY(name, endpoint_id)
		);
	`
	_, err = pool.Exec(context.Background(), queryParameterTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating query parameter table", err)
		os.Exit(1)
	}

	// Create request body fields table
	requestBodyFieldsTableQuery := `
		CREATE TABLE IF NOT EXISTS request_body_fields (
			name VARCHAR NOT NULL,
			type VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			required BOOLEAN NOT NULL,
			endpoint_id VARCHAR,
			CONSTRAINT request_body_field_endpoint_id_fkey FOREIGN KEY(endpoint_id) REFERENCES endpoints(endpoint_id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			PRIMARY KEY(name, endpoint_id)
		);
	`
	_, err = pool.Exec(context.Background(), requestBodyFieldsTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request body fields table", err)
		os.Exit(1)
	}

	// Create headers table
	headersTableQuery := `
		CREATE TABLE IF NOT EXISTS headers (
			name VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			required BOOLEAN NOT NULL,
			endpoint_id VARCHAR,
			CONSTRAINT header_endpoint_id_fkey FOREIGN KEY(endpoint_id) REFERENCES endpoints(endpoint_id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			PRIMARY KEY(name, endpoint_id)
		);
	`
	_, err = pool.Exec(context.Background(), headersTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating headers table", err)
		os.Exit(1)
	}

	// Create status code table
	statusCodeTableQuery := `
		CREATE TABLE IF NOT EXISTS status_codes (
			code INTEGER NOT NULL,
			description VARCHAR NOT NULL,
			endpoint_id VARCHAR,
			CONSTRAINT status_code_endpoint_id_fkey FOREIGN KEY(endpoint_id) REFERENCES endpoints(endpoint_id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			PRIMARY KEY(code, endpoint_id)
		);
	`
	_, err = pool.Exec(context.Background(), statusCodeTableQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating status codes table", err)
		os.Exit(1)
	}
}
