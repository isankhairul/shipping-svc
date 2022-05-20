package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"math"
	"strings"

	"gorm.io/gorm"
)

type courierRepo struct {
	base BaseRepository
}

type CourierRepository interface {
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Courier, *base.Pagination, error)
	FindByUid(uid *string) (*entity.Courier, error)
	CreateCourier(product *entity.Courier) (*entity.Courier, error)
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB
	Delete(uid string) error
	Update(uid string, input map[string]interface{}) error
}

func NewCourierRepository(br BaseRepository) CourierRepository {
	return &courierRepo{br}
}

func (r *courierRepo) FindByUid(uid *string) (*entity.Courier, error) {
	var courier entity.Courier
	err := r.base.GetDB().
		Where("uid=?", uid).
		First(&courier).Error
	if err != nil {
		return nil, err
	}

	return &courier, nil
}

func (r *courierRepo) CreateCourier(courier *entity.Courier) (*entity.Courier, error) {
	err := r.base.GetDB().
		Create(courier).Error
	if err != nil {
		return nil, err
	}

	return courier, nil
}

func (r *courierRepo) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(currRecord)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *courierRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Courier, *base.Pagination, error) {
	var couriers []entity.Courier
	var pagination base.Pagination

	query := r.base.GetDB()

	if filter["courier_type"] != "" {
		query = query.Where("LOWER(courier_type) LIKE ?", "%"+strings.ToLower(filter["courier_type"].(string))+"%")
	}

	if filter["status"] != "" {
		query = query.Where("LOWER(status) LIKE ?", "%"+strings.ToLower(filter["status"].(string))+"%")
	}

	if len(sort) > 0 {
		query = query.Order(sort)
	}

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(couriers, &pagination, query, int64(len(couriers)))).
		Find(&couriers).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return couriers, &pagination, nil
}

func (r *courierRepo) Delete(uid string) error {
	var courier entity.Courier
	err := r.base.GetDB().
		Where("uid = ?", uid).
		Delete(&courier).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r *courierRepo) Update(uid string, input map[string]interface{}) error {
	err := r.base.GetDB().Model(&entity.Courier{}).
		Where("uid=?", uid).
		Updates(input).Error
	if err != nil {
		return err
	}
	return nil
}