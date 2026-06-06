package repository

import (
	"org_struct_api/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *EmployeeRepository) FindByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.First(&employee, id).Error; err != nil {
		return nil, mapError(err)
	}
	return &employee, nil
}

func (r *EmployeeRepository) FindByDepartmentID(departmentID uint) ([]models.Employee, error) {
	var employees []models.Employee
	query := r.db.Where("department_id = ?", departmentID).Order("full_name ASC")
	if err := query.Find(&employees).Error; err != nil {
		return nil, mapError(err)
	}
	return employees, nil
}

func (r *EmployeeRepository) ReassignDepartment(fromDepartmentID uint, toDepartmentID uint) error {
	return r.db.Model(&models.Employee{}).Where("department_id = ?", fromDepartmentID).Update("department_id", toDepartmentID).Error
}
