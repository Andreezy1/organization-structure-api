package models

import "time"

type Employee struct {
	ID           uint       `gorm:"primaryKey"`
	DepartmentID uint       `gorm:"not null"`
	FullName     string     `gorm:"size:200;not null"`
	Position     string     `gorm:"size:200;not null"`
	HiredAt      *time.Time `gorm:"type:date"`
	CreatedAt    time.Time  `gorm:"not null;default:now()"`
}
