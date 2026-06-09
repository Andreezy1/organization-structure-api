package handler

import (
	"encoding/json"
	"net/http"
	"org_struct_api/internal/models"
)

type DepartmentService interface {
	CreateDepartment(department *models.Department) (*models.Department, error)
	DeleteDepartment(departmentID uint, mode string, reassignTo *uint) error
	GetDepartment(id uint, depth int, includeEmployees bool) (*models.Department, error)
	UpdateDepartment(id uint, name string, parentID models.Patch[*uint]) (*models.Department, error)
}

type DepartmentHandler struct {
	departmentService DepartmentService
}

func NewDepartmentHandler(departmentService DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
	}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var department CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	createdDepartment, err := h.departmentService.CreateDepartment(&models.Department{Name: department.Name,
		ParentID: department.ParentID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, createdDepartment)
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		writeError(w, http.StatusBadRequest, "mode is required")
		return
	}

	reassignTo, err := queryUintPtr(r, "reassign_to_department_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.departmentService.DeleteDepartment(uint(departmentID), mode, reassignTo)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	var departmentRequest UpdateDepartmentRequest
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedDepartment, err := h.departmentService.UpdateDepartment(uint(id), departmentRequest.Name, departmentRequest.ParentID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedDepartment)
}

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	includeEmployees := true

	if value := r.URL.Query().Get("include_employees"); value != "" {
		includeEmployees = value == "true"
	}

	depth, err := queryIntDefault(r, "depth", 1)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	department, err := h.departmentService.GetDepartment(uint(id), depth, includeEmployees)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, department)
}
