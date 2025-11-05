package handler

import (
	"encoding/json"
	"net/http"

	"github.com/example/task-management/application/command"
	"github.com/example/task-management/application/dto"
	"github.com/example/task-management/application/query"
	"github.com/example/task-management/interface/http/middleware"
	"github.com/example/task-management/shared/di"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	container    *di.Container
	errorHandler *middleware.ErrorHandler
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(container *di.Container) *TaskHandler {
	return &TaskHandler{
		container:    container,
		errorHandler: middleware.NewErrorHandler(),
	}
}

// CreateTask handles POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create command
	cmd := command.CreateTaskCommand{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		AssigneeID:  req.AssigneeID,
		Deadline:    req.Deadline,
		CreatedBy:   r.Header.Get("X-User-ID"), // In real app, from auth context
	}

	// Handle command
	result, err := h.container.CreateTaskCommandHandler.Handle(cmd)
	if err != nil {
		httpErr := h.errorHandler.HandleError(err)
		h.writeJSON(w, httpErr.Code, httpErr)
		return
	}

	// Return response
	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"task_id": result.TaskID,
		"message": "Task created successfully",
	})
}

// GetTask handles GET /tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		h.writeError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	// Create query
	q := query.GetTaskQuery{
		TaskID: taskID,
	}

	// Handle query
	result, err := h.container.GetTaskQueryHandler.Handle(q)
	if err != nil {
		httpErr := h.errorHandler.HandleError(err)
		h.writeJSON(w, httpErr.Code, httpErr)
		return
	}

	// Return response
	h.writeJSON(w, http.StatusOK, result)
}

// ListTasksByProject handles GET /projects/{id}/tasks
func (h *TaskHandler) ListTasksByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		h.writeError(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	status := r.URL.Query().Get("status")

	// Create query
	q := query.ListTasksByProjectQuery{
		ProjectID: projectID,
		Status:    status,
	}

	// Handle query
	results, err := h.container.ListTasksByProjectQueryHandler.Handle(q)
	if err != nil {
		httpErr := h.errorHandler.HandleError(err)
		h.writeJSON(w, httpErr.Code, httpErr)
		return
	}

	// Return response
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"tasks": results,
		"count": len(results),
	})
}

// AssignTask handles POST /tasks/{id}/assign
func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		h.writeError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	var req dto.AssignTaskRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create command
	cmd := command.AssignTaskCommand{
		TaskID:     taskID,
		AssigneeID: req.AssigneeID,
		AssignedBy: r.Header.Get("X-User-ID"),
	}

	// Handle command
	_, err := h.container.AssignTaskCommandHandler.Handle(cmd)
	if err != nil {
		httpErr := h.errorHandler.HandleError(err)
		h.writeJSON(w, httpErr.Code, httpErr)
		return
	}

	// Return response
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Task assigned successfully",
	})
}

// UpdateTaskStatus handles PUT /tasks/{id}/status
func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		h.writeError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	var req dto.UpdateTaskStatusRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create command
	cmd := command.UpdateTaskStatusCommand{
		TaskID:    taskID,
		NewStatus: req.Status,
	}

	// Handle command
	_, err := h.container.UpdateTaskStatusCommandHandler.Handle(cmd)
	if err != nil {
		httpErr := h.errorHandler.HandleError(err)
		h.writeJSON(w, httpErr.Code, httpErr)
		return
	}

	// Return response
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Task status updated successfully",
	})
}

// Helper methods

// writeJSON writes a JSON response
func (h *TaskHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func (h *TaskHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSON(w, statusCode, map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}