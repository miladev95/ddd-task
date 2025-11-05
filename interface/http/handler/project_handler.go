package handler

import (
	"encoding/json"
	"net/http"

	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
	"github.com/miladev95/ddd-task/interface/http/middleware"
	"github.com/miladev95/ddd-task/shared/di"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	container    *di.Container
	errorHandler *middleware.ErrorHandler
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(container *di.Container) *ProjectHandler {
	return &ProjectHandler{
		container:    container,
		errorHandler: middleware.NewErrorHandler(),
	}
}

// CreateProjectRequest represents the request to create a project
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id" binding:"required"`
	WorkflowID  string `json:"workflow_id" binding:"required"`
}

// CreateProject handles POST /api/projects
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate owner exists
	ownerID, err := value.NewUserID(req.OwnerID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}

	_, err = h.container.UserRepository.GetByID(ownerID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Owner user not found")
		return
	}

	// Validate workflow exists
	workflowID, err := value.NewWorkflowID(req.WorkflowID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid workflow ID")
		return
	}

	_, err = h.container.WorkflowRepository.GetByID(workflowID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Workflow not found")
		return
	}

	// Generate new project ID
	projectID := value.GenerateProjectID()

	// Create project aggregate
	project, err := aggregate.NewProject(projectID, req.Name, req.Description, ownerID, workflowID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Save project
	err = h.container.ProjectRepository.Save(project)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to save project")
		return
	}

	// Return response
	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"project_id": projectID.Value(),
		"name":       project.Name(),
		"message":    "Project created successfully",
	})
}

// GetProject handles GET /api/projects/{id}
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("id")
	if projectID == "" {
		h.writeError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	// Create project ID
	id, err := value.NewProjectID(projectID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get project
	project, err := h.container.ProjectRepository.GetByID(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Project not found")
		return
	}

	// Return response
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":          project.ID().Value(),
		"name":        project.Name(),
		"description": project.Description(),
		"owner_id":    project.OwnerID().Value(),
		"workflow_id": project.WorkflowID().Value(),
		"task_count":  project.TaskCount(),
		"created_at":  project.CreatedAt(),
		"updated_at":  project.UpdatedAt(),
		"archived":    project.IsArchived(),
	})
}

// Helper methods

// writeJSON writes a JSON response
func (h *ProjectHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func (h *ProjectHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSON(w, statusCode, map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}