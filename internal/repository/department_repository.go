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

// func (dr *DepartmentRepository) FindByNameAndParent(
// 	name string,
// 	parentID *uint,
// ) (*models.Department, error) {
// 	var department models.Department
// 	query := dr.db.Where("name = ?", name)
// 	if parentID == nil {
// 		query = query.Where("parent_id IS NULL")
// 	} else {
// 		query = query.Where("parent_id = ?", *parentID)
// 	}
// 	err := query.First(&department).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &department, nil
// }

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

func (r *DepartmentRepository) GetRootDepartments() ([]models.Department, error) {
	var departments []models.Department
	query := r.db.Preload("Children").Where("parent_id IS NULL")
	err := query.Find(&departments).Error
	if err != nil {
		return nil, mapError(err)
	}
	return departments, nil
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

func (r *DepartmentRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Department{}, id)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mapError(result.Error)
	}
	return nil
}

func (r *DepartmentRepository) Update(department *models.Department) error {
	result := r.db.Save(department)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mapError(result.Error)
	}
	return nil
}
