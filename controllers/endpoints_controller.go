package controllers

import (
	"slices"
	"context"
	"net/http"
	"crypto/rand"
	"encoding/json"
	"encoding/base64"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/calvintran1478/api-doc-server/templ"
	"github.com/calvintran1478/api-doc-server/utils"
)

type EndpointsController struct {
	Pool *pgxpool.Pool
}

type AddEndpointRequest struct {
	Method 		string `json:"method"`
	Path 		string `json:"path"`
	Description string `json:"description"`
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

	// Send success response
	endpoint := templ.Endpoint{EndpointID: endpointID, Method: body.Method, Path: body.Path, Description: body.Description}
	w.WriteHeader(http.StatusCreated)
	templ.EndpointEntry(endpoint).Render(r.Context(), w)
}
