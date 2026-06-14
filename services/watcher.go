package services

import (
	"dev-audit-tracker/models"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// Fungsi utama yang di panggil di main.go
func WatchWorkspace(db *gorm.DB, workspacePath string) {
	log.Printf("Watching workspace: %s", workspacePath)
	entries, _ := os.ReadDir(workspacePath)
	for _, entry := range entries {
		if entry.IsDir() {
			ticketID := entry.Name()
			log.Printf("Processing ticket: %s", ticketID)

			var project models.Project
			log.Printf("Processing ticket: %s", ticketID)
			// Cari project berdasarkan ticket_number
			result := db.Where("ticket_number = ?", ticketID).First(&project)
			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					// Jika tidak ada, buat baru
					project = models.Project{
						TicketNumber: ticketID,
						Status:       "CREATED",
					}
					if err := db.Create(&project).Error; err != nil {
						log.Printf("Error creating project for %s: %v", ticketID, err)
						continue
					}
					log.Printf("Created new project for %s with ID %d", ticketID, project.ID)
				} else {
					log.Printf("Error finding project for %s: %v", ticketID, result.Error)
					continue
				}
			} else {
				log.Printf("Found existing project for %s with ID %d", ticketID, project.ID)
			}
			// Pastikan ID valid
			if project.ID == 0 {
				log.Printf("Warning: Project ID is 0 for %s, skipping", ticketID)
				continue
			}

			// Jalankan logika pengecekan sub-folder
			ticketPath := filepath.Join(workspacePath, ticketID)
			log.Printf("ticket path: %s", ticketPath)
			checkSubFolder(db, &project, ticketPath)
			log.Printf("Finished processing ticket: %s with status: %s", ticketID, project.Status)
		}
	}
}

func checkSubFolder(db *gorm.DB, p *models.Project, ticketPath string) {
	log.Printf("Checking subfolders for ticket: %s", p.TicketNumber)
	// cek FSD avalilable
	fsdFiles, _ := filepath.Glob(filepath.Join(ticketPath, "FSD", "*.pdf"))
	log.Printf("Found %d FSD files for ticket: %s", len(fsdFiles), p.TicketNumber)
	if len(fsdFiles) > 0 && !p.HasFSD {
		p.HasFSD = true
		UpdateStatus(db, p, "KICK_OFF", "File FSD terdeteksi")
	}

	// 2. Cek Analysis (Status: DEVELOPMENT_START)
	analysisFiles, _ := filepath.Glob(filepath.Join(ticketPath, "Analysis", "*.xls*"))
	if len(analysisFiles) > 0 && !p.HasAnalysis && p.HasFSD {
		p.HasAnalysis = true
		UpdateStatus(db, p, "DEVELOPMENT_START", "File Analisis terdeteksi")
	}

	// 3. Cek SIT (Status: SIT_DONE)
	sitFiles, _ := filepath.Glob(filepath.Join(ticketPath, "SIT", "*SIT_RESULT_DONE*.xls*"))

	if len(sitFiles) > 0 && !p.HasSIT {
		p.HasSIT = true
		log.Printf("Updating status to SIT_DONE for ticket: %s", p.TicketNumber)
		UpdateStatus(db, p, "SIT_DONE", "File SIT Result Done terdeteksi")
	}

	// 3. Cek Document (Status: DEV_DONE jika ada 4 file DONE)
	docFiles, _ := filepath.Glob(filepath.Join(ticketPath, "Document", "*DONE*"))
	if len(docFiles) >= 4 && p.DocCount < 4 && p.HasAnalysis {
		p.DocCount = len(docFiles)
		UpdateStatus(db, p, "DEV_DONE", "4 Dokumen DONE terdeteksi")
	}

	// 4. Cek Revisi (Status: DEV_REVISION_UAT)
	revisiFile := filepath.Join(ticketPath, "Revisi", "revisi.txt")
	if _, err := os.Stat(revisiFile); err == nil && p.Status != "DEV_REVISION_UAT" {
		UpdateStatus(db, p, "DEV_REVISION_UAT", "File revisi.txt ditemukan")
	}

	// 5. Cek PROJECT DONE (Status: PROJECT_DONE)
	ccbFile := filepath.Join(ticketPath, "CCB", "ccb.txt")
	if _, err := os.Stat(ccbFile); err == nil && p.Status != "PROJECT_DONE" {
		UpdateStatus(db, p, "PROJECT_DONE", "File ccb.txt ditemukan")
	}
}

func UpdateStatus(db *gorm.DB, p *models.Project, newStatus string, note string) {
	if p.ID == 0 {
		log.Printf("Error: Cannot update status for project with ID 0")
		return
	}
	log.Printf("Updating status for project %s to %s: %s", p.TicketNumber, newStatus, note)
	p.Status = newStatus
	p.UpdatedAt = time.Now()
	if err := db.Save(p).Error; err != nil {
		log.Printf("Error saving project %s: %v", p.TicketNumber, err)
		return
	}
	if err := db.Create(&models.ProjectHistory{
		ProjectID: p.ID,
		Status:    newStatus,
		Notes:     note,
		CreatedAt: time.Now(),
	}).Error; err != nil {
		log.Printf("Error creating history for project %s: %v", p.TicketNumber, err)
	}
}
