package models

import "time"

type Department struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name string `gorm:"size:200;not null" json:"name"`

	ParentID *uint `json:"parent_id"`

	Parent *Department `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	Children []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	Employees []Employee `json:"employees,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
