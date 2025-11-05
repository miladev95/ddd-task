package handler

import (
	"encoding/json"
	"net/http"

	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
	"github.com/miladev95/ddd-task/interface/http/middleware"
	"github.com/miladev95/ddd-task/shared/di"
)

// WorkflowHandler handles HTTP requests for workflows
type WorkflowHandler struct {
	container    *di.Container
	errorHandler *middleware.ErrorHandler
}

// NewWorkflowHandler creates a new WorkflowHandler
func NewWorkflowHandler(container *di.Container) *WorkflowHandler {
	return &WorkflowHandler{
		container:    container,
		errorHandler: middleware.NewErrorHandler(),
	}
}

// WorkflowStatusRequest represents a workflow status in the request
type WorkflowStatusRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order" binding:"required"`
	IsFinal     bool   `json:"is_final"`
}

// CreateWorkflowRequest represents the request to create a workflow
type CreateWorkflowRequest struct {
	Name        string                    `json:"name" binding:"required"`
	Description string                    `json:"description"`
	Statuses    []WorkflowStatusRequest   `json:"statuses" binding:"required"`
}

// CreateWorkflow handles POST /api/workflows
func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkflowRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate at least one status
	if len(req.Statuses) == 0 {
		h.writeError(w, http.StatusBadRequest, "At least one status is required")
		return
	}

	// Convert statuses
	statuses := make([]aggregate.WorkflowStatus, len(req.Statuses))
	for i, s := range req.Statuses {
		statuses[i] = aggregate.NewWorkflowStatus(s.Name, s.Description, s.Order, s.IsFinal)
	}

	// Generate new workflow ID
	workflowID := value.GenerateWorkflowID()

	// Create workflow aggregate
	workflow, err := aggregate.NewWorkflow(workflowID, req.Name, req.Description, statuses)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Save workflow
	err = h.container.WorkflowRepository.Save(workflow)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to save workflow")
		return
	}

	// Return response
	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"workflow_id": workflowID.Value(),
		"name":        workflow.Name(),
		"description": workflow.Description(),
		"message":     "Workflow created successfully",
	})
}

// GetWorkflow handles GET /api/workflows/{id}
func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := r.URL.Query().Get("id")
	if workflowID == "" {
		h.writeError(w, http.StatusBadRequest, "Workflow ID is required")
		return
	}

	// Create workflow ID
	id, err := value.NewWorkflowID(workflowID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid workflow ID")
		return
	}

	// Get workflow
	workflow, err := h.container.WorkflowRepository.GetByID(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Workflow not found")
		return
	}

	// Return response
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":          workflow.ID().Value(),
		"name":        workflow.Name(),
		"description": workflow.Description(),
		"created_at":  workflow.CreatedAt(),
		"updated_at":  workflow.UpdatedAt(),
	})
}

// Helper methods

// writeJSON writes a JSON response
func (h *WorkflowHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func (h *WorkflowHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSON(w, statusCode, map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}