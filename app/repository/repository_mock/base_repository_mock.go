package repository_mock

import (
	"go-klikdokter/app/model/base"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type BaseRepositoryMock struct {
	Mock mock.Mock
}

func (b *BaseRepositoryMock) GetDB() *gorm.DB {
	return nil
}

func (b *BaseRepositoryMock) BeginTx() {
	/*
		implemented
	*/
}

func (b *BaseRepositoryMock) CommitTx() {
	/*
		implemented
	*/
}

func (b *BaseRepositoryMock) RollbackTx() {
	/*
		implemented
	*/
}

func (b *BaseRepositoryMock) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	return nil
}
