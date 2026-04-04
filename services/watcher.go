package services

import (
	"dev-audit-tracker/models"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// Fungsi utama yang di panggil di main.go
func WatchWorkspace(db *gorm.DB, workspacePath string) {
	entries, _ := os.ReadDir(workspacePath)
	for _, entry := range entries {
		if entry.IsDir() {
			ticketID := entry.Name()

			var project models.Project
			// Cari atau buat data proejct di DB
			db.Where("ticket_number = ?", ticketID).FirstOrCreate(&project, models.Project{
				TicketNumber: ticketID,
				Status:       "CREATED",
			})

			// Jalankan logika pengecekan sub-folder
			ticketPath := filepath.Join(workspacePath, ticketID)
			checkSubFolder(db, &project, ticketPath)
		}
	}
}

func checkSubFolder(db *gorm.DB, p *models.Project, ticketPath string) {
	// cek FSD avalilable
	fsdFiles, _ := filepath.Glob(filepath.Join(ticketPath, "FSD", "*.pdf"))
	if len(fsdFiles) > 0 && !p.HasFSD {
		p.HasFSD = true
		UpdateStatus(db, p, "KICK_OFF", "File FSD terdeteksi")
	}
}

func UpdateStatus(db *gorm.DB, p *models.Project, newStatus string, note string) {
	p.Status = newStatus
	p.UpdatedAt = time.Now()
	db.Save(p)
	db.Create(&models.ProjectHistory{
		ProjectID: p.ID,
		Status:    newStatus,
		Notes:     note,
		CreatedAt: time.Now(),
	})
}
