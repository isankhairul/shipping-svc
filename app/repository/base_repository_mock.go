package repository

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

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

}

func (b *BaseRepositoryMock) CommitTx() {

}

func (b *BaseRepositoryMock) RollbackTx() {

}

func (b *BaseRepositoryMock) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	return nil
}

func (b *BaseRepositoryMock) FindByUid(uid string, model interface{}) (interface{}, error) {
	arguments := b.Mock.Called(uid)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		product := arguments.Get(0).(entity.Product)
		return &product, nil
	}
}

func (b *BaseRepositoryMock) FindById(id string, entity interface{}) (interface{}, error) {
	arguments := b.Mock.Called(id)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		return &entity, nil
	}
}

func (b *BaseRepositoryMock) Create(input interface{}) (interface{}, error) {
	return input, nil
}

func (b *BaseRepositoryMock) UpdateByUid(uid string, input map[string]interface{}, entity interface{}) error {
	return nil
}

func (b *BaseRepositoryMock) UpdateById(id string, input map[string]interface{}, entity interface{}) error {
	return nil
}

func (b *BaseRepositoryMock) DeleteByUid(uid string, entity interface{}) error {
	return nil
}

func (b *BaseRepositoryMock) DeleteById(id string, entity interface{}) error {
	return nil
}
