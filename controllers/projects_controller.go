package controllers

import (
	"io"
	"context"
	"errors"
	"unicode"
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/calvintran1478/api-doc-server/templ"
	"github.com/calvintran1478/api-doc-server/utils"
)

type ProjectsController struct {
	Pool *pgxpool.Pool
}

/*
 * Adds a project for the user
 *
 * Method: POST
 * Path: /api/v1/projects
 */
func (c *ProjectsController) AddProject(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}

	// Read request body
	name, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	for _, runeValue := range string(name) {
		if !unicode.IsLetter(runeValue) && !unicode.IsDigit(runeValue) && !unicode.IsSpace(runeValue) {
			http.Error(w, "Name can only contain letters, numbers, and spaces", http.StatusBadRequest)
			return
		}
	}

	// Generate project ID
	projectIDBytes := make([]byte, 16)
	_, err = rand.Read(projectIDBytes)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	projectID := base64.RawURLEncoding.EncodeToString(projectIDBytes)

	// Add project
	commandTag, err := c.Pool.Exec(context.Background(), "INSERT INTO projects (user_id, project_id, name) VALUES ($1, $2, $3) ON CONFLICT (name) DO NOTHING", userID, projectID, name)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if commandTag.RowsAffected() != 1 {
		http.Error(w, "Project with name already exists", http.StatusConflict)
		return
	}

	// Send success response
	project := templ.Project{ProjectID: projectID, Name: string(name)}
	w.WriteHeader(http.StatusCreated)
	templ.ProjectEntry(project).Render(r.Context(), w)
}

/*
 * Gets HTML page of user projects
 *
 * Method: GET
 * Path: /projects
 */
func (c *ProjectsController) GetProjects(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}

	// Fetch projects
	rows, err := c.Pool.Query(context.Background(), "SELECT project_id, name FROM projects WHERE user_id=$1", userID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var projects []templ.Project
	for rows.Next() {
		var project templ.Project
		err = rows.Scan(&project.ProjectID, &project.Name)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		projects = append(projects, project)
	}

	// Render HTML from projects
	templ.ProjectsPage(projects).Render(r.Context(), w)
}

/*
 * Gets HTML page of a single user project
 *
 * Method: GET
 * Path: /projects/{projectID}
 */
func (c *ProjectsController) GetProject(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}
	projectID := r.PathValue("projectID")

	// Fetch project
	var project templ.Project
	err := c.Pool.QueryRow(context.Background(), "SELECT project_id, name FROM projects WHERE user_id=$1 AND project_id=$2", userID, projectID).Scan(&project.ProjectID, &project.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Render HTML from project
	templ.ProjectPage(project).Render(r.Context(), w)
}

/*
 * Gets HTML page of project settings for a single user project
 *
 * Method: GET
 * Path: /projects/{projectID}/settings
 */
func (c *ProjectsController) GetProjectSettings(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}
	projectID := r.PathValue("projectID")

	// Fetch project
	var project templ.Project
	err := c.Pool.QueryRow(context.Background(), "SELECT name FROM projects WHERE user_id=$1 AND project_id=$2", userID, projectID).Scan(&project.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	project.ProjectID = projectID

	// Render HTML from project
	templ.SettingsPage(project).Render(r.Context(), w)
}

/*
 * Updates the name of a user project
 *
 * Method: PATCH
 * Path: /api/projects/{projectID}
 */
func (c *ProjectsController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}
	projectID := r.PathValue("projectID")

	// Read request body
	name, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	for _, runeValue := range string(name) {
		if !unicode.IsLetter(runeValue) && !unicode.IsDigit(runeValue) && !unicode.IsSpace(runeValue) {
			http.Error(w, "Name can only contain letters, numbers, and spaces", http.StatusBadRequest)
			return
		}
	}

	// Check if name already exists
	var nameExists bool
	err = c.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM projects WHERE user_id=$1 AND name=$2 AND project_id<>$3)", userID, name, projectID).Scan(&nameExists)
	if nameExists {
		http.Error(w, "Project with the given name already exists", http.StatusConflict)
		return
	}

	// Update project
	commandTag, err := c.Pool.Exec(context.Background(), "UPDATE projects SET name=$3 WHERE user_id=$1 AND projects_id=$2", userID, projectID, name)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if commandTag.RowsAffected() != 1 {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusNoContent)
}

/*
 * Deletes a user project
 *
 * Method: DELETE
 * Path: /api/projects/{projectID}
 */
func (c *ProjectsController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}
	projectID := r.PathValue("projectID")

	// Delete project
	commandTag, err := c.Pool.Exec(context.Background(), "DELETE FROM projects WHERE user_id=$1 AND project_id=$2", userID, projectID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if commandTag.RowsAffected() != 1 {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusNoContent)
}
