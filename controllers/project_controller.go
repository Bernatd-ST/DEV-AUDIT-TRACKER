package controllers

import (
	"dev-audit-tracker/models"
	"net/http"
	"text/template"

	"gorm.io/gorm"
)

type ProjectController struct {
	DB *gorm.DB
}

func (pc *ProjectController) Dashboard(w http.ResponseWriter, r *http.Request) {
	var projects []models.Project
	pc.DB.Order("updated_at desc").Find(&projects)

	tmpl, err := template.ParseFiles("./views/index.html")
	if err != nil {
		http.Error(w, "View tidak ditemukan", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Project": projects,
	}

	tmpl.Execute(w, data)
}
