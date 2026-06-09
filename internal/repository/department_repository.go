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

func (r *DepartmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *DepartmentRepository) FindByID(id uint) (*models.Department, error) {
	var department models.Department
	err := r.db.First(&department, id).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &department, nil
}

func (r *DepartmentRepository) ExistsByNameAndParent(
	name string,
	parentID *uint,
) (bool, error) {
	var count int64
	query := r.db.
		Model(&models.Department{}).
		Where("name = ?", name)
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, mapError(err)
	}
	return count > 0, nil
}

func (r *DepartmentRepository) ExistsByID(id uint) (bool, error) {
	var count int64
	err := r.db.
		Model(&models.Department{}).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		return false, mapError(err)
	}
	return count > 0, nil
}

func (r *DepartmentRepository) GetChildren(parentID uint) ([]models.Department, error) {
	var departments []models.Department
	query := r.db.Where("parent_id = ?", parentID)
	err := query.Find(&departments).Error
	if err != nil {
		return nil, mapError(err)
	}
	return departments, nil
}

func (r *DepartmentRepository) GetChildrenBatch(parentIDs []uint) ([]models.Department, error) {
	var departments []models.Department
	err := r.db.
		Where("parent_id IN ?", parentIDs).
		Find(&departments).Error
	if err != nil {
		return nil, mapError(err)
	}
	return departments, nil
}

func (r *DepartmentRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Department{}, id)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *DepartmentRepository) Update(department *models.Department) error {
	result := r.db.Save(department)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}
