package controllers

import (
	"dev-audit-tracker/models"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gorm.io/gorm"
)

type ProjectController struct {
	DB *gorm.DB
}

func (pc *ProjectController) Dashboard(w http.ResponseWriter, r *http.Request) {
	var projects []models.Project
	pc.DB.Order("id").Find(&projects)

	// Prepare view objects with URL-escaped ticket (to be used in links)
	type ProjectView struct {
		Project       models.Project
		TicketEscaped string
	}
	var views []ProjectView
	for _, p := range projects {
		views = append(views, ProjectView{Project: p, TicketEscaped: url.QueryEscape(p.TicketNumber)})
	}

	tmpl, err := template.ParseFiles("./views/index.html")
	if err != nil {
		http.Error(w, "View tidak ditemukan", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Projects": views,
	}

	tmpl.Execute(w, data)
}

func (pc *ProjectController) LogFrontend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&logData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	message, ok := logData["message"].(string)
	if !ok {
		message = "Unknown log message"
	}

	log.Printf("[FRONTEND] %s", message)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged"))
}

func (pc *ProjectController) Comparison(w http.ResponseWriter, r *http.Request) {
	// Support query param `?ticket=...` or path /comparison/{ticket}
	ticketID := r.URL.Query().Get("ticket")
	if ticketID == "" {
		ticketID = strings.TrimPrefix(r.URL.Path, "/comparison/")
	}
	// Unescape if needed
	if u, err := url.QueryUnescape(ticketID); err == nil && u != "" {
		ticketID = u
	}
	var project models.Project
	result := pc.DB.Where("ticket_number = ?", ticketID).First(&project)
	if result.Error != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	// read comparison folder
	// adjust workspace path if needed
	workspacePath := "/Users/bernatdsitumeang/Desktop/Ticketing"
	compPath := filepath.Join(workspacePath, ticketID, "comparison")
	files := []string{}
	if entries, err := os.ReadDir(compPath); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				files = append(files, e.Name())
			}
		}
	}

	// render comparison view
	tmpl, err := template.ParseFiles("./views/comparison.html")
	if err != nil {
		http.Error(w, "View tidak ditemukan", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Project": project,
		"Files":   files,
	}

	tmpl.Execute(w, data)
}

func (pc *ProjectController) Projects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var projects []models.Project
	pc.DB.Order("id").Find(&projects)
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, "Failed to encode projects", http.StatusInternalServerError)
		return
	}
}
