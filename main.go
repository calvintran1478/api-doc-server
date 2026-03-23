package main

import (
	"os"
	"fmt"
	"context"
	"strconv"
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/calvintran1478/api-doc-server/controllers"
	"github.com/calvintran1478/api-doc-server/utils"
)

func main() {
	// Load enviornment variables
	utils.LoadEnv(".env")
	bcryptCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid bcrypt cost", err)
		os.Exit(1)
	} else if bcryptCost <= 0 {
		fmt.Fprintf(os.Stderr, "Bcrypt cost should be positive", err)
		os.Exit(1)
	}

	// Set up database connection
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_CONN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Set up resource controllers
	userController := controllers.UserController{Pool: pool, BcryptCost: bcryptCost}
	projectsController := controllers.ProjectsController{Pool: pool}

	// Set up file server
	fileServer := http.FileServer(http.Dir("./static/"))

	// Set up routes
	router := http.NewServeMux()
	router.HandleFunc("POST /api/register", userController.RegisterUser)
	router.HandleFunc("POST /api/login", userController.LoginUser)
	router.HandleFunc("POST /api/logout", userController.LogoutUser)
	router.HandleFunc("POST /api/projects", projectsController.AddProject)
	router.HandleFunc("GET /projects", projectsController.GetProjects)
	router.HandleFunc("GET /projects/{projectID}", projectsController.GetProject)
	router.HandleFunc("PATCH /api/projects/{projectID}", projectsController.UpdateProject)
	router.HandleFunc("DELETE /api/projects/{projectID}", projectsController.DeleteProject)
	router.HandleFunc("GET /projects/{projectID}/settings", projectsController.GetProjectSettings)
	router.Handle("/", fileServer)

	// Initialize and start up server
	server := http.Server{
		Addr:	 ":8080",
		Handler: router,
	}

	fmt.Println("Listening on http://localhost:8080")
	server.ListenAndServe()
}
