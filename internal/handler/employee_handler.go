package handler

import (
	"encoding/json"
	"net/http"
	"org_struct_api/internal/models"
)

type EmployeeService interface {
	CreateEmployee(employee *models.Employee) (*models.Employee, error)
}

type EmployeeHandler struct {
	employeeService EmployeeService
}

func NewEmployeeHandler(employeeService EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee CreateEmployeeRequest

	id, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdEmployee, err := h.employeeService.CreateEmployee(&models.Employee{
		DepartmentID: id,
		FullName:     employee.FullName,
		Position:     employee.Position,
		HiredAt:      employee.HiredAt,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdEmployee)
}
