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
	// Initialize handlers
	taskHandler := handler.NewTaskHandler(r.container)
	projectHandler := handler.NewProjectHandler(r.container)
	userHandler := handler.NewUserHandler(r.container)
	workflowHandler := handler.NewWorkflowHandler(r.container)

	// User routes
	r.mux.HandleFunc("/api/users", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			userHandler.CreateUser(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/users/get", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			userHandler.GetUser(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Workflow routes
	r.mux.HandleFunc("/api/workflows", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			workflowHandler.CreateWorkflow(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/workflows/get", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			workflowHandler.GetWorkflow(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Project routes
	r.mux.HandleFunc("/api/projects", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			projectHandler.CreateProject(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/projects/get", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			projectHandler.GetProject(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Task routes
	r.mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			taskHandler.CreateTask(w, req)
		case http.MethodGet:
			taskHandler.ListTasksByProject(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/tasks/get", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			taskHandler.GetTask(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/tasks/assign", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			taskHandler.AssignTask(w, req)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/api/tasks/status", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			taskHandler.UpdateTaskStatus(w, req)
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