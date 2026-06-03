package repository

import (
	"org_struct_api/internal/models"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		db: db,
	}
}

func (dr *DepartmentRepository) Create(department *models.Department) error {
	return dr.db.Create(department).Error
}

func (dr *DepartmentRepository) FindByID(id uint) (*models.Department, error) {
	var department models.Department

	err := dr.db.First(&department, id).Error

	if err != nil {
		return nil, err
	}

	return &department, nil
}

func (dr *DepartmentRepository) FindByNameAndParent(
	name string,
	parentID *uint,
) (*models.Department, error) {

	var department models.Department

	query := dr.db.Where("name = ?", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.First(&department).Error

	if err != nil {
		return nil, err
	}

	return &department, nil
}

func (dr *DepartmentRepository) GetRootDepartments() ([]models.Department, error) {
	var departments []models.Department

	query := dr.db.Preload("Children").Where("parent_id IS NULL")

	err := query.Find(&departments).Error

	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (dr *DepartmentRepository) GetChildren(parentID uint) ([]models.Department, error) {
	var departments []models.Department

	query := dr.db.Where("parent_id = ?", parentID)

	err := query.Find(&departments).Error

	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (dr *DepartmentRepository) Delete(id uint) error {
	return dr.db.Delete(&models.Department{}, id).Error
}

func (dr *DepartmentRepository) Update(department *models.Department) error {
	return dr.db.Save(department).Error
}
