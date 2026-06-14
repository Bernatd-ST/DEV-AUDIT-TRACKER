package routes

import (
	"dev-audit-tracker/controllers"
	"net/http"

	"gorm.io/gorm"
)

func enableCORS(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		enableCORS(w, req)
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, req)
	}
}

func InitRoutes(db *gorm.DB) {
	projectController := &controllers.ProjectController{DB: db}

	// Routes untuk dashboard
	http.HandleFunc("/", withCORS(projectController.Dashboard))
	// Route untuk frontend logging
	http.HandleFunc("/log", withCORS(projectController.LogFrontend))
	// route untuk button project done
	http.HandleFunc("/comparison/", withCORS(projectController.Comparison))

	http.HandleFunc("/projects", withCORS(projectController.Projects))
}
