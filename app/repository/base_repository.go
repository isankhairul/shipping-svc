package repository

import (
	"math"

	"go-klikdokter/app/model/base"

	"gorm.io/gorm"
)

type baseRepository struct {
	db *gorm.DB
}

type BaseRepository interface {
	GetDB() *gorm.DB
	BeginTx()
	CommitTx()
	RollbackTx()
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB
	FindByUid(uid string, model interface{}) (interface{}, error)
	FindById(id string, entity interface{}) (interface{}, error)
	Create(input interface{}) (interface{}, error)
	UpdateByUid(uid string, input map[string]interface{}, entity interface{}) error
	UpdateById(id string, input map[string]interface{}, entity interface{}) error
	DeleteByUid(uid string, entity interface{}) error
	DeleteById(id string, entity interface{}) error
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &baseRepository{db}
}

func (br *baseRepository) GetDB() *gorm.DB {
	return br.db
}

func (br *baseRepository) BeginTx() {
	br.db = br.GetDB().Begin()
}

func (br *baseRepository) CommitTx() {
	br.GetDB().Commit()
}

func (br *baseRepository) RollbackTx() {
	br.GetDB().Rollback()
}

func (br *baseRepository) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(currRecord)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

//Generic standard CRUD function
func (br *baseRepository) FindByUid(uid string, entity interface{}) (interface{}, error) {
	err := br.GetDB().
		Where("uid=?", uid).
		First(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (br *baseRepository) FindById(id string, entity interface{}) (interface{}, error) {
	err := br.GetDB().
		Where("id=?", id).
		First(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (br *baseRepository) Create(input interface{}) (interface{}, error) {
	err := br.GetDB().Create(input).Error
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (br *baseRepository) UpdateByUid(uid string, input map[string]interface{}, entity interface{}) error {
	err := br.GetDB().Model(entity).
		Where("uid=?", uid).
		Updates(input).Error
	if err != nil {
		return err
	}

	return nil
}

func (br *baseRepository) UpdateById(id string, input map[string]interface{}, entity interface{}) error {
	err := br.GetDB().Model(entity).
		Where("id=?", id).
		Updates(input).Error
	if err != nil {
		return err
	}

	return nil
}

func (br *baseRepository) DeleteByUid(uid string, entity interface{}) error {
	err := br.GetDB().
		Where("uid = ?", uid).
		Delete(entity).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (br *baseRepository) DeleteById(id string, entity interface{}) error {
	err := br.GetDB().
		Where("id = ?", id).
		Delete(entity).
		Error
	if err != nil {
		return err
	}

	return nil
}
