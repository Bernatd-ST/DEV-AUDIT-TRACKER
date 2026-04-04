package routes

import (
	"dev-audit-tracker/controllers"
	"net/http"

	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB) {
	projectController := &controllers.ProjectController{DB: db}

	// Routes untuk dashboard
	http.HandleFunc("/", projectController.Dashboard)
}
