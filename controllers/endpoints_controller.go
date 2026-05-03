package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/calvintran1478/api-doc-server/templ"
	"github.com/calvintran1478/api-doc-server/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"slices"
	"sort"
)

type EndpointsController struct {
	Pool *pgxpool.Pool
}

type PathParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type AddEndpointRequest struct {
	Method         string          `json:"method"`
	Path           string          `json:"path"`
	Description    string          `json:"description"`
	PathParameters []PathParameter `json:"path_parameters"`
}

/*
 * Adds an endpoint for the user
 *
 * Method: POST
 * Path: /api/projects/{projectID}/endpoints
 */
func (c *EndpointsController) AddEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}
	projectID := r.PathValue("projectID")

	// Read request body
	var body AddEndpointRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
		return
	}

	// Validate request body
	httpMethods := []string{"POST", "GET", "PATCH", "DELETE", "PUT"}
	if !slices.Contains(httpMethods, body.Method) {
		http.Error(w, "Invalid method for endpoint", http.StatusBadRequest)
		return
	}
	sort.Slice(body.PathParameters, func(i, j int) bool {
		return body.PathParameters[i].Name < body.PathParameters[j].Name
	})
	for i := range len(body.PathParameters) - 1 {
		if body.PathParameters[i].Name == body.PathParameters[i+1].Name {
			http.Error(w, "Duplicate path parameter found", http.StatusConflict)
			return
		}
	}

	// Generate endpoint ID
	endpointIDBytes := make([]byte, 16)
	_, err = rand.Read(endpointIDBytes)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	endpointID := base64.RawURLEncoding.EncodeToString(endpointIDBytes)

	// Add endpoint
	commandTag, err := c.Pool.Exec(context.Background(), "INSERT INTO endpoints (project_id, endpoint_id, method, path, description) SELECT $1, $2, $3, $4, $5 WHERE EXISTS (SELECT 1 FROM projects WHERE project_id=$6)", projectID, endpointID, body.Method, body.Path, body.Description, projectID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if commandTag.RowsAffected() != 1 {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Add path parameters
	for _, pathParameter := range body.PathParameters {
		commandTag, err = c.Pool.Exec(context.Background(), "INSERT INTO path_parameters (endpoint_id, name, type, description) VALUES ($1, $2, $3, $4)", endpointID, pathParameter.Name, pathParameter.Type, pathParameter.Description)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	// Send success response
	endpoint := templ.Endpoint{EndpointID: endpointID, Method: body.Method, Path: body.Path, Description: body.Description}
	w.WriteHeader(http.StatusCreated)
	templ.EndpointEntry(endpoint).Render(r.Context(), w)
}

/*
 * Deletes an endpoint from a user project
 *
 * Method: DELETE
 * Path: /api/projects/{projectID}/endpoints/{endpointID}
 */
func (c *EndpointsController) DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID := utils.GetUser(w, r)
	if userID == "" {
		return
	}
	projectID := r.PathValue("projectID")
	endpointID := r.PathValue("endpointID")

	// Delete endpoint
	commandTag, err := c.Pool.Exec(context.Background(), "DELETE FROM endpoints WHERE project_id=$1 AND endpoint_id=$2 AND EXISTS (SELECT 1 FROM projects WHERE user_id=$3 AND project_id=$4)", projectID, endpointID, userID, projectID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if commandTag.RowsAffected() != 1 {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusNoContent)
}
