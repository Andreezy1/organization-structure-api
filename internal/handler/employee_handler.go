package handler

import (
	"encoding/json"
	"net/http"
	"org_struct_api/internal/dto"
	"org_struct_api/internal/models"
	"org_struct_api/internal/service"
)

type EmployeeService interface {
	CreateEmployee(employee *models.Employee) (*models.Employee, error)
	GetDepartmentEmployees(departmentID uint) ([]models.Employee, error)
}

type EmployeeHandler struct {
	employeeService EmployeeService
}

func NewEmployeeHandler(employeeService *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee dto.CreateEmployeeRequest

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

func (h *EmployeeHandler) GetDepartmentEmployees(w http.ResponseWriter, r *http.Request) {
	departmentID, err := parsePathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	employees, err := h.employeeService.GetDepartmentEmployees(uint(departmentID))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, employees)
}
