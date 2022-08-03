package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"math"

	"gorm.io/gorm"
)

type courierServiceRepo struct {
	base BaseRepository
}

type CourierServiceRepository interface {
	FindAll(limit int, page int, sort string) ([]entity.CourierService, *base.Pagination, error)
	CheckExistsByCourierIdShippingCode(courierUId string, shippingCode string) (bool, error)
	CheckExistsByUIdCourierIdShippingCode(uid string, courierUId string, shippingCode string) (bool, error)
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.CourierService, *base.Pagination, error)
	FindByUid(uid *string) (*entity.CourierService, error)
	CreateCourierService(courierservice *entity.CourierService) (*entity.CourierService, error)
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB
	Delete(uid string) error
	Update(uid string, input map[string]interface{}) error
	IsCourierServiceAssigned(courierServiceID uint64) bool
}

func NewCourierServiceRepository(br BaseRepository) CourierServiceRepository {
	return &courierServiceRepo{br}
}

func (r *courierServiceRepo) FindByUid(uid *string) (*entity.CourierService, error) {
	var courierService entity.CourierService
	err := r.base.GetDB().Preload("Courier").
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

func (r *courierServiceRepo) FindAll(limit int, page int, sort string) ([]entity.CourierService, *base.Pagination, error) {
	var courierService []entity.CourierService
	var pagination base.Pagination

	query := r.base.GetDB()
	if len(sort) > 0 {
		query = query.Order(sort)
	} else {
		query = query.Order("updated_at DESC")
	}
	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(courierService, &pagination, query, int64(len(courierService)))).
		Find(&courierService).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return courierService, &pagination, nil
}

func (r *courierServiceRepo) CheckExistsByCourierIdShippingCode(courierUId string, shippingCode string) (bool, error) {
	var exists bool
	err := r.base.GetDB().
		Model(&entity.CourierService{}).
		Select("count(*) > 0").
		Where("shipping_code = ? AND courier_uid = ?", shippingCode, courierUId).
		Find(&exists).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r *courierServiceRepo) CheckExistsByUIdCourierIdShippingCode(uid string, courierUId string, shippingCode string) (bool, error) {
	var exists bool
	err := r.base.GetDB().
		Model(&entity.CourierService{}).
		Select("count(*) > 0").
		Where("uid != ? AND shipping_code = ? AND courier_uid = ?", uid, shippingCode, courierUId).
		Find(&exists).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (r *courierServiceRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.CourierService, *base.Pagination, error) {
	var couriers []entity.CourierService
	var pagination base.Pagination

	query := r.base.GetDB().Model(&entity.CourierService{}).Preload("Courier").Joins("Courier")

	if filter["courier_uid"] != "" {
		query = query.Where("courier_uid = ?", filter["courier_uid"])
	}

	if filter["shipping_code"] != "" {
		query = query.Where("shipping_code = ?", filter["shipping_code"])
	}

	if filter["shipping_name"] != "" {
		query = query.Where("shipping_name = ?", filter["shipping_name"])
	}

	if filter["status"] != 0 {
		//query = query.Where("status = ?", filter["status"])
		status := filter["status"].(int32)
		query = query.Where(&entity.CourierService{Status: &status})
	}

	if len(sort) > 0 {
		query = query.Order(sort)
	} else {
		query = query.Order("courier_service.updated_at DESC")
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

func (r *courierServiceRepo) IsCourierServiceAssigned(courierServiceID uint64) bool {
	var count int64
	r.base.GetDB().Model(&entity.ChannelCourierService{}).
		Where(&entity.ChannelCourierService{CourierServiceID: courierServiceID}).
		Count(&count)

	return count > 0
}
