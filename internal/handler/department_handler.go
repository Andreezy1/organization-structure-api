package handler

import (
	"encoding/json"
	"net/http"
	"org_struct_api/internal/dto"
	"org_struct_api/internal/models"
	"strconv"
)

type DepartmentService interface {
	CreateDepartment(department *models.Department) (*models.Department, error)
	DeleteDepartment(departmentID uint, mode string, reassignTo *uint) error
	GetDepartment(id uint, depth int, includeEmployees bool) (*models.Department, error)
	GetDepartmentTree() ([]models.Department, error)
	UpdateDepartment(id uint, name string, parentID *uint) (*models.Department, error)
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

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var department models.Department

	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "error: invalid request body"})
		return
	}
	createdDepartment, err := h.departmentService.CreateDepartment(&department)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDepartment)

}

func (h *DepartmentHandler) GetDepartmentTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	departments, err := h.departmentService.GetDepartmentTree()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(departments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	departmentIDstr := r.PathValue("id")
	mode := r.URL.Query().Get("mode")
	reassign_to_department_id_str := r.URL.Query().Get("reassign_to_department_id")

	if departmentIDstr == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "department id is required",
		})

		return
	}

	departmentID, err := strconv.Atoi(departmentIDstr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	if mode == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "mode required",
		})

		return
	}
	var reassignTo *uint

	if reassign_to_department_id_str != "" {
		id, err := strconv.Atoi(reassign_to_department_id_str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

			return
		}
		tmp := uint(id)
		reassignTo = &tmp
	}

	err = h.departmentService.DeleteDepartment(uint(departmentID), mode, reassignTo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "department deleted",
	})
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	var departmentRequest dto.UpdateDepartmentRequest

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "department id is required",
		})

		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "error: invalid request body"})
		return
	}

	updatedDepartment, err := h.departmentService.UpdateDepartment(uint(id), departmentRequest.Name, departmentRequest.ParentID)
	if err != nil {
		switch err.Error() {

		case "department not found":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusConflict)
		}

		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedDepartment)

}

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "department id is required",
		})

		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}

	includeEmployees := true

	if value := r.URL.Query().Get("include_employees"); value != "" {
		includeEmployees = value == "true"
	}

	depth := 1

	depthStr := r.URL.Query().Get("depth")

	if depthStr != "" {
		depth, err = strconv.Atoi(depthStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

			return
		}
	}

	department, err := h.departmentService.GetDepartment(uint(id), depth, includeEmployees)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(department)

}
