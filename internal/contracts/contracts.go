package contracts

import "org_struct_api/internal/models"

type DeportmentRepo interface {
	Create(department *models.Department) error
	Delete(id uint) error
	ExistsByID(id uint) (bool, error)
	ExistsByNameAndParent(name string, parentID *uint) (bool, error)
	FindByID(id uint) (*models.Department, error)
	GetChildren(parentID uint) ([]models.Department, error)
	Update(department *models.Department) error
	GetChildrenBatch(parentIDs []uint) ([]models.Department, error)
}

type TxManager interface {
	InTx(fn func(DeportmentRepo, EmployeeRepo) error) error
}

type EmployeeRepo interface {
	Create(employee *models.Employee) error
	FindByDepartmentID(departmentID uint) ([]models.Employee, error)
	FindByID(id uint) (*models.Employee, error)
	ReassignDepartment(fromDepartmentID uint, toDepartmentID uint) error
}
