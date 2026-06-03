package models

import "time"

type Employee struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DepartmentID uint `gorm:"not null" json:"department_id"`

	Department Department `gorm:"foreignKey:DepartmentID" json:"department"`

	FullName string `gorm:"size:200;not null" json:"full_name"`

	Position string `gorm:"size:200;not null" json:"position"`

	HiredAt *time.Time `json:"hired_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
