package service

import (
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
		return nil, err
	}
	position, err := validate("position", employee.Position)
	if err != nil {
		return nil, err
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
