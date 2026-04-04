package main

import (
	"dev-audit-tracker/models"
	"dev-audit-tracker/routes"
	"dev-audit-tracker/services"
	"fmt"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Inisialisasi Database
	dsn := "root:@tcp(127.0.0.1:3307)/dev_audit_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 2. Auto Migrate
	db.AutoMigrate(&models.Project{}, &models.ProjectHistory{})
	fmt.Println("✅ Database Connected & Migrated")

	// 3. Setup Workspace
	workspacePath := "C:\\Users\\user\\Desktop\\workspace-Bernatdev\\"
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		os.Mkdir(workspacePath, 0755)
	}

	// 4. Jalankan Background Watcher (Goroutine)
	go func() {
		fmt.Println("🔍 Service Watcher Running...")
		for {
			services.WatchWorkspace(db, workspacePath)
			time.Sleep(5 * time.Second)
		}
	}()

	// 5. Inisialisasi Routes & Server
	routes.InitRoutes(db)

	fmt.Println("🌍 Dashboard: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
