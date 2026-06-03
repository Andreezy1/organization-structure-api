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

func (er *EmployeeRepository) Create(employee *models.Employee) error {
	return er.db.Create(employee).Error
}

func (er *EmployeeRepository) FindByID(id uint) (*models.Employee, error) {
	var employee models.Employee

	if err := er.db.First(&employee, id).Error; err != nil {
		return nil, err
	}

	return &employee, nil
}

func (er *EmployeeRepository) FindByDepartmentID(departmentID uint) ([]models.Employee, error) {
	var employees []models.Employee
	query := er.db.Where("department_id = ?", departmentID).Order("full_name ASC")

	if err := query.Find(&employees).Error; err != nil {
		return nil, err
	}

	return employees, nil
}

func (er *EmployeeRepository) ReassignDepartment(fromDepartmentID uint, toDepartmentID uint) error {
	return er.db.Model(&models.Employee{}).Where("department_id = ?", fromDepartmentID).Update("department_id", toDepartmentID).Error
}
