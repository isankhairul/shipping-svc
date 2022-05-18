package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"math"
	"strings"

	"gorm.io/gorm"
)

type courierServiceRepo struct {
	base BaseRepository
}

type CourierServiceRepository interface {
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.CourierService, *base.Pagination, error)
	FindByUid(uid *string) (*entity.CourierService, error)
	CreateCourierService(product *entity.CourierService) (*entity.CourierService, error)
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB
	Delete(uid string) error
	Update(uid string, input map[string]interface{}) error
}

func NewCourierServiceRepository(br BaseRepository) CourierServiceRepository {
	return &courierServiceRepo{br}
}

func (r *courierServiceRepo) FindByUid(uid *string) (*entity.CourierService, error) {
	var courierService entity.CourierService
	err := r.base.GetDB().
		Where("uid=?", uid).
		First(&courierService).Error
	if err != nil {
		return nil, err
	}

	return &courierService, nil
}

func (r *courierServiceRepo) CreateCourierService(courierService *entity.CourierService) (*entity.CourierService, error) {
	err := r.base.GetDB().
		Create(courierService).Error
	if err != nil {
		return nil, err
	}

	return courierService, nil
}

func (r *courierServiceRepo) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(currRecord)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *courierServiceRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.CourierService, *base.Pagination, error) {
	var couriers []entity.CourierService
	var pagination base.Pagination

	query := r.base.GetDB()

	if filter["shipping_type"] != "" {
		query = query.Where("LOWER(shipping_type) LIKE ?", "%"+strings.ToLower(filter["shipping_type"].(string))+"%")
	}

	if filter["status"] != "" {
		query = query.Where("status = ?", filter["status"])
	}

	if filter["courier_id"] != "" {
		query = query.Where("courier_id = ?", filter["courier_id"])
	}

	if filter["shipping_code"] != "" {
		query = query.Where("shipping_code = ?", filter["shipping_code"])
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

func (r *courierServiceRepo) Delete(uid string) error {
	var courierService entity.CourierService
	err := r.base.GetDB().
		Where("uid = ?", uid).
		Delete(&courierService).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r *courierServiceRepo) Update(uid string, input map[string]interface{}) error {
	err := r.base.GetDB().Model(&entity.CourierService{}).
		Where("uid=?", uid).
		Updates(input).Error
	if err != nil {
		return err
	}
	return nil
}
