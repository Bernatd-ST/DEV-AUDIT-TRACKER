package main

import (
	"dev-audit-tracker/models"
	"dev-audit-tracker/routes"
	"dev-audit-tracker/services"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Setup logging to file
	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	// 1. Inisialisasi Database
	dsn := "root:@tcp(127.0.0.1:3307)/dev_audit_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 2. Auto Migrate
	db.AutoMigrate(&models.Project{}, &models.ProjectHistory{})
	log.Println("✅ Database Connected & Migrated")

	// 3. Setup Workspace
	workspacePath := "/Users/bernatdsitumeang/Desktop/Ticketing/"
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		os.Mkdir(workspacePath, 0755)
	}

	// 4. Jalankan Background Watcher (Goroutine)
	go func() {
		log.Println("🔍 Service Watcher Running...")
		for {
			services.WatchWorkspace(db, workspacePath)
			time.Sleep(5 * time.Second)
		}
	}()

	// 5. Inisialisasi Routes & Server
	routes.InitRoutes(db)

	log.Println("🌍 Dashboard: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
