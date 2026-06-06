package models

import "time"

type Department struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null"`
	ParentID  *uint
	CreatedAt time.Time    `gorm:"not null;default:now()"`
	Children  []Department `gorm:"foreignKey:ParentID"`
	Employees []Employee
}
