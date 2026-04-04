package models

import "time"

type Project struct {
	ID           uint   `gorm:"primaryKey"`
	TicketNumber string `gorm:"unique;not null"`
	Status       string `gorm:"default:'CREATED'"`
	HasFSD       bool   `gorm:"default:false"`
	HasAnalysis  bool   `gorm:"default:false"`
	DocCount     int    `gorm:"default:0"` // untuk jumlah 4 dokumen signature
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Histories    []ProjectHistory
}

type ProjectHistory struct {
	ID        uint `gorm:"primaryKey"`
	ProjectID uint
	Status    string
	Notes     string
	CreatedAt time.Time
}
