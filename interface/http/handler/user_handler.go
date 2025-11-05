package handler

import (
	"encoding/json"
	"net/http"

	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
	"github.com/miladev95/ddd-task/interface/http/middleware"
	"github.com/miladev95/ddd-task/shared/di"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	container    *di.Container
	errorHandler *middleware.ErrorHandler
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(container *di.Container) *UserHandler {
	return &UserHandler{
		container:    container,
		errorHandler: middleware.NewErrorHandler(),
	}
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate new user ID
	userID := value.GenerateUserID()

	// Create user aggregate
	user, err := aggregate.NewUser(userID, req.Email, req.FirstName, req.LastName)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Save user
	err = h.container.UserRepository.Save(user)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to save user")
		return
	}

	// Return response
	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"user_id":    userID.Value(),
		"email":      user.Email(),
		"first_name": user.FirstName(),
		"last_name":  user.LastName(),
		"full_name":  user.FullName(),
		"message":    "User created successfully",
	})
}

// GetUser handles GET /api/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		h.writeError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Create user ID
	id, err := value.NewUserID(userID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get user
	user, err := h.container.UserRepository.GetByID(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "User not found")
		return
	}

	// Return response
	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":         user.ID().Value(),
		"email":      user.Email(),
		"first_name": user.FirstName(),
		"last_name":  user.LastName(),
		"full_name":  user.FullName(),
		"created_at": user.CreatedAt(),
		"updated_at": user.UpdatedAt(),
	})
}

// Helper methods

// writeJSON writes a JSON response
func (h *UserHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func (h *UserHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSON(w, statusCode, map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}