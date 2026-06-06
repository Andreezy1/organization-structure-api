package repository

import (
	"org_struct_api/internal/contracts"

	"gorm.io/gorm"
)

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) InTx(fn func(contracts.DeportmentRepo, contracts.EmployeeRepo) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		return fn(&DepartmentRepository{db: tx}, &EmployeeRepository{db: tx})
	})
}
