package controllers

import (
	"io"
	"context"
	"bytes"
	"errors"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

type UserController struct {
	Pool *pgxpool.Pool
	BcryptCost int
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	args := bytes.Split(body, []byte("\n"))

	// Get request body fields
	email := args[0]
	password := args[1]
	firstName := args[2]
	lastName := args[3]

	// Check if a user with the given email already exists
	var exists bool
	err = c.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "User with email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), c.BcryptCost)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Register user in the database
	_, err = c.Pool.Exec(context.Background(), "INSERT INTO users (user_id, email, password, first_name, last_name) VALUES ($1, $2, $3, $4, $5)", uuid.New(), email, hashedPassword, firstName, lastName)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
}

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	args := bytes.Split(body, []byte("\n"))

	// Get request body fields
	email := args[0]
	password := args[1]

	// Look up user in database
	var userID string
	var hashedPassword string
	err = c.Pool.QueryRow(context.Background(), "SELECT user_id, password FROM users WHERE email=$1", email).Scan(&userID, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "User with email not found", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Verify user password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Add access token
	cookie := http.Cookie{
		Name: "access-token",
		Value: userID,
		MaxAge: 604800,
		HttpOnly: true,
		Secure: true,
		Path: "/",
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	// Send success response
	w.WriteHeader(http.StatusOK)
}
