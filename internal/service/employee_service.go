package service

import (
	"errors"
	"org_struct_api/internal/models"
	"org_struct_api/internal/repository"
	"strings"
)

type EmployeeService struct {
	employeeRepo   *repository.EmployeeRepository
	departmentRepo *repository.DepartmentRepository
}

func NewEmployeeService(employeeRepo *repository.EmployeeRepository,
	departmentRepo *repository.DepartmentRepository) *EmployeeService {
	return &EmployeeService{employeeRepo: employeeRepo,
		departmentRepo: departmentRepo}
}

func (es *EmployeeService) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	employee.FullName = strings.TrimSpace(employee.FullName)
	employee.Position = strings.TrimSpace(employee.Position)

	if employee.FullName == "" {
		return nil, errors.New("employee name is required")
	}

	if employee.Position == "" {
		return nil, errors.New("employee position is required")
	}

	_, err := es.departmentRepo.FindByID(employee.DepartmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	err = es.employeeRepo.Create(employee)

	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (es *EmployeeService) GetDepartmentEmployees(departmentID uint) ([]models.Employee, error) {
	_, err := es.departmentRepo.FindByID(departmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}
	employees, err := es.employeeRepo.FindByDepartmentID(departmentID)
	if err != nil {
		return nil, err
	}
	return employees, nil
}
