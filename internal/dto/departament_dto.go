package dto

import "org_struct_api/internal/models"

type UpdateDepartmentRequest struct {
	Name     string              `json:"name"`
	ParentID models.Patch[*uint] `json:"parent_id"`
}

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}
