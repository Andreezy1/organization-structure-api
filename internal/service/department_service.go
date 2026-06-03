package service

import (
	"errors"
	"org_struct_api/internal/models"
	"org_struct_api/internal/repository"
	"strings"

	"gorm.io/gorm"
)

type DepartmentService struct {
	departmentRepo *repository.DepartmentRepository
	employeeRepo   *repository.EmployeeRepository
}

func NewDepartmentService(departmentRepo *repository.DepartmentRepository,
	employeeRepo *repository.EmployeeRepository) *DepartmentService {
	return &DepartmentService{
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
	}
}

func (ds *DepartmentService) CreateDepartment(department *models.Department) (*models.Department, error) {
	department.Name = strings.TrimSpace(department.Name)

	if department.Name == "" {
		return nil, errors.New("department name is required")
	}

	if department.ParentID != nil {
		_, err := ds.departmentRepo.FindByID(*department.ParentID)
		if err != nil {
			return nil, errors.New("parent department not found")
		}
	}

	existingDepartment, err := ds.departmentRepo.FindByNameAndParent(department.Name, department.ParentID)
	if err == nil && existingDepartment != nil {
		return nil, errors.New("department alredy exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	err = ds.departmentRepo.Create(department)

	if err != nil {
		return nil, err
	}
	return department, nil
}

func (ds *DepartmentService) GetDepartmentTree() ([]models.Department, error) {

	departments, err := ds.departmentRepo.GetRootDepartments()

	if err != nil {
		return nil, err
	}

	for i := range departments {
		err := ds.loadChildren(&departments[i], 100)
		if err != nil {
			return nil, err
		}
	}

	return departments, nil
}

func (ds *DepartmentService) loadChildren(department *models.Department, depth int) error {
	if depth <= 0 {
		return nil
	}

	children, err := ds.departmentRepo.GetChildren(department.ID)

	if err != nil {
		return err
	}

	department.Children = children

	for i := range department.Children {

		err := ds.loadChildren(&department.Children[i], depth-1)

		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *DepartmentService) DeleteDepartment(departmentID uint, mode string, reassignTo *uint) error {
	_, err := ds.departmentRepo.FindByID(departmentID)
	if err != nil {
		return err
	}
	switch {
	case mode == "cascade":
		return ds.departmentRepo.Delete(departmentID)
	case mode == "reassign":
		if reassignTo == nil {
			return errors.New("reassign_to_department_id is required")
		}
		_, err := ds.departmentRepo.FindByID(*reassignTo)
		if err != nil {
			return err
		}
		if *reassignTo == departmentID {
			return errors.New("reassign_to_department_id = department_id")
		}
		err = ds.employeeRepo.ReassignDepartment(departmentID, *reassignTo)
		if err != nil {
			return err
		}
		return ds.departmentRepo.Delete(departmentID)
	default:
		return errors.New("invalide mode")
	}
}

func (ds *DepartmentService) isDescendant(departmentID uint, targetID uint) (bool, error) {
	children, err := ds.departmentRepo.GetChildren(departmentID)
	if err != nil {
		return false, err
	}
	for _, child := range children {
		if child.ID == targetID {
			return true, nil
		}
		found, err := ds.isDescendant(child.ID, targetID)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

func (ds *DepartmentService) UpdateDepartment(id uint, name string, parentID *uint) (*models.Department, error) {
	department, err := ds.departmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	if name != "" {
		department.Name = name
	}
	if parentID != nil && *parentID == id {
		return nil, errors.New("department cannot be parent of itself")
	}
	if parentID != nil {
		_, err := ds.departmentRepo.FindByID(*parentID)
		if err != nil {
			return nil, err
		}
	}

	if parentID != nil {
		isCycle, err := ds.isDescendant(id, *parentID)

		if err != nil {
			return nil, err
		}

		if isCycle {
			return nil, errors.New("cyclic hierarchy detected")
		}
	}

	newParentID := department.ParentID

	if parentID != nil {
		newParentID = parentID
	}

	existingDepartment, err := ds.departmentRepo.FindByNameAndParent(department.Name, newParentID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err == nil && existingDepartment != nil && existingDepartment.ID != department.ID {
		return nil, errors.New("department already exists")
	}

	department.ParentID = newParentID

	err = ds.departmentRepo.Update(department)

	if err != nil {
		return nil, err
	}

	return department, nil
}

func (ds *DepartmentService) GetDepartment(id uint, depth int, includeEmployees bool) (*models.Department, error) {
	department, err := ds.departmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}

	err = ds.loadChildren(department, depth)

	if err != nil {
		return nil, err
	}

	if includeEmployees {
		employees, err := ds.employeeRepo.FindByDepartmentID(id)

		if err != nil {
			return nil, err
		}
		department.Employees = employees
	}

	return department, nil
}
