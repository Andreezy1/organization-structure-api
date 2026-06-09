package models

import "time"

type Department struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Name      string        `gorm:"size:200;not null" json:"name"`
	ParentID  *uint         `json:"parent_id"`
	CreatedAt time.Time     `gorm:"not null;default:now()" json:"created_at"`
	Children  []*Department `gorm:"foreignKey:ParentID" json:"children"`
	Employees []Employee    `json:"employees"`
}
