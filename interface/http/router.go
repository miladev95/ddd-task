package http

import (
	"net/http"

	"github.com/miladev95/ddd-task/interface/http/handler"
	"github.com/miladev95/ddd-task/shared/di"
)

// Router sets up all HTTP routes
type Router struct {
	container    *di.Container
	mux          *http.ServeMux
	taskHandler  *handler.TaskHandler
}

// NewRouter creates a new Router
func NewRouter(container *di.Container) *Router {
	return &Router{
		container:   container,
		mux:         http.NewServeMux(),
		taskHandler: handler.NewTaskHandler(container),
	}
}

// SetupRoutes sets up all HTTP routes
func (r *Router) SetupRoutes() {
	// Task routes
	r.mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.taskHandler.CreateTask(w, req)
		case http.MethodGet:
			r.taskHandler.ListTasksByProject(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/tasks/get", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			r.taskHandler.GetTask(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/tasks/assign", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			r.taskHandler.AssignTask(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/tasks/status", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			r.taskHandler.UpdateTaskStatus(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	r.mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})
}

// Handler returns the HTTP handler
func (r *Router) Handler() http.Handler {
	return r.mux
}