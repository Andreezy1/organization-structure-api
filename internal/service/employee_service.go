package service

import (
	"errors"
	"fmt"
	"org_struct_api/internal/contracts"
	"org_struct_api/internal/models"
)

type EmployeeService struct {
	employeeRepo   contracts.EmployeeRepo
	departmentRepo contracts.DeportmentRepo
}

func NewEmployeeService(employeeRepo contracts.EmployeeRepo,
	departmentRepo contracts.DeportmentRepo) *EmployeeService {
	return &EmployeeService{employeeRepo: employeeRepo,
		departmentRepo: departmentRepo}
}

func (s *EmployeeService) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	fullname, err := validate("name", employee.FullName)
	if err != nil {
		return nil, models.ErrValidation
	}
	position, err := validate("name", employee.Position)
	if err != nil {
		return nil, models.ErrValidation
	}
	exists, err := s.departmentRepo.ExistsByID(employee.DepartmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("%w: department", models.ErrNotFound)
	}

	employee.FullName = fullname
	employee.Position = position
	err = s.employeeRepo.Create(employee)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *EmployeeService) GetDepartmentEmployees(departmentID uint) ([]models.Employee, error) {
	_, err := s.departmentRepo.FindByID(departmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}
	employees, err := s.employeeRepo.FindByDepartmentID(departmentID)
	if err != nil {
		return nil, err
	}
	return employees, nil
}
