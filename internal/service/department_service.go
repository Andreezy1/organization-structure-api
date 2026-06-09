package service

import (
	"fmt"

	"org_struct_api/internal/contracts"
	"org_struct_api/internal/models"
)

const (
	maxTreeDepth     = 5
	defaultTreeDepth = 1
)

type DepartmentService struct {
	departmentRepo contracts.DeportmentRepo
	employeeRepo   contracts.EmployeeRepo
	txManager      contracts.TxManager
}

func NewDepartmentService(
	departmentRepo contracts.DeportmentRepo,
	employeeRepo contracts.EmployeeRepo,
	txManager contracts.TxManager,
) *DepartmentService {
	return &DepartmentService{
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
		txManager:      txManager,
	}
}

func (s *DepartmentService) CreateDepartment(department *models.Department) (*models.Department, error) {
	name, err := validate("name", department.Name)
	if err != nil {
		return nil, err
	}

	if department.ParentID != nil {
		exists, err := s.departmentRepo.ExistsByID(*department.ParentID)
		if err != nil {
			return nil, fmt.Errorf("check parent: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("%w: parent department", models.ErrNotFound)
		}
	}

	existingDepartment, err := s.departmentRepo.ExistsByNameAndParent(department.Name, department.ParentID)
	if err != nil {
		return nil, fmt.Errorf("check parent: %w", err)
	}
	if existingDepartment {
		return nil, fmt.Errorf("%w: already exists", models.ErrConflict)
	}

	department.Name = name
	err = s.departmentRepo.Create(department)
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (s *DepartmentService) DeleteDepartment(departmentID uint, mode string, reassignTo *uint) error {
	switch mode {
	case "cascade":
		return s.departmentRepo.Delete(departmentID)
	case "reassign":
		if reassignTo == nil {
			return fmt.Errorf("%w: reassign_to_department_id is required", models.ErrValidation)
		}
		exists, err := s.departmentRepo.ExistsByID(*reassignTo)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("%w: reassign_to_department_id", models.ErrNotFound)
		}
		if *reassignTo == departmentID {
			return fmt.Errorf("%w,: reassign_to_department_id = department_id", models.ErrValidation)
		}
		return s.txManager.InTx(func(dr contracts.DeportmentRepo, er contracts.EmployeeRepo) error {
			err := er.ReassignDepartment(departmentID, *reassignTo)
			if err != nil {
				return err
			}
			if err := dr.Delete(departmentID); err != nil {
				return err
			}
			return nil
		})
	default:
		return fmt.Errorf("%w: invalide mode", models.ErrValidation)
	}
}

func (s *DepartmentService) isDescendant(departmentID uint, targetID uint) (bool, error) {
	children, err := s.departmentRepo.GetChildren(departmentID)
	if err != nil {
		return false, err
	}
	for _, child := range children {
		if child.ID == targetID {
			return true, nil
		}
		found, err := s.isDescendant(child.ID, targetID)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

func (s *DepartmentService) UpdateDepartment(id uint, name string, parentID models.Patch[*uint]) (*models.Department, error) {
	department, err := s.departmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		name, err = validate("name", name)
		if err != nil {
			return nil, err
		}
		department.Name = name
	}

	newParentID := department.ParentID
	if parentID.Set {
		newParentID = parentID.Value
	}
	if newParentID != nil {
		if *newParentID == department.ID {
			return nil, fmt.Errorf("%w: department cannot be parent of itself", models.ErrValidation)
		}
		exist, err := s.departmentRepo.ExistsByID(*newParentID)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, fmt.Errorf("%w: new parent department", models.ErrNotFound)
		}
		isCycle, err := s.isDescendant(id, *newParentID)
		if err != nil {
			return nil, err
		}
		if isCycle {
			return nil, models.ErrCycle
		}
	}

	conflict, err := s.departmentRepo.ExistsByNameAndParent(department.Name, newParentID)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, fmt.Errorf("%w: department already exists", models.ErrConflict)
	}

	department.ParentID = newParentID
	err = s.departmentRepo.Update(department)
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (s *DepartmentService) GetDepartment(id uint, depth int, includeEmployees bool) (*models.Department, error) {
	department, err := s.departmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if depth < 1 {
		depth = defaultTreeDepth
	}
	if depth > maxTreeDepth {
		depth = maxTreeDepth
	}

	err = s.loadChildrenBFS(department, depth)
	if err != nil {
		return nil, err
	}

	if includeEmployees {
		employees, err := s.employeeRepo.FindByDepartmentID(id)
		if err != nil {
			return nil, err
		}
		department.Employees = employees
	}
	return department, nil
}

func (s *DepartmentService) loadChildrenBFS(root *models.Department, maxDepth int) error {
	nodes := map[uint]*models.Department{
		root.ID: root,
	}
	currentLevel := []uint{root.ID}
	for level := 0; level < maxDepth; level++ {
		children, err := s.departmentRepo.GetChildrenBatch(currentLevel)
		if err != nil {
			return err
		}
		if len(children) == 0 {
			break
		}
		nextLevel := make([]uint, 0, len(children))
		for i := range children {
			child := &children[i]
			nodes[child.ID] = child
			nextLevel = append(nextLevel, child.ID)
		}
		for i := range children {
			child := &children[i]
			if child.ParentID == nil {
				continue
			}
			parent := nodes[*child.ParentID]
			parent.Children = append(parent.Children, child)
		}
		currentLevel = nextLevel
	}
	return nil
}
