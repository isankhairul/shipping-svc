package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"
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
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]response.CourierServiceListResponse, *base.Pagination, error)
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

func (r *courierServiceRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]response.CourierServiceListResponse, *base.Pagination, error) {
	var couriers []entity.CourierService
	var pagination base.Pagination
	var couriersResponse []response.CourierServiceListResponse

	//query := r.base.GetDB().Model(&entity.CourierService{}).Preload("Courier").Joins("Courier")

	query := r.base.GetDB().Model(&entity.CourierService{}).
		Select("courier_service.*, \"Courier\".courier_name, \"Courier\".courier_type,sp1.title as courier_type_name, sp2.title as shipping_type_name").
		//Preload("Courier").
		Joins("Courier").
		Joins("LEFT JOIN shippment_predefined sp1 ON sp1.code  = \"Courier\".courier_type ").
		Joins("LEFT JOIN shippment_predefined sp2 ON sp2.code  = courier_service.shipping_type ").
		Where("sp1.type = 'courier_type'").
		Where("sp2.type = 'shipping_type'")

	for k, v := range filter {
		switch k {
		case "shipping_code", "shipping_name":
			value, ok := v.([]string)
			if ok && len(value) > 0 {
				query = query.Where(like(k, value))

			}
		case "courier_uid", "courier_type", "shipping_type":
			value, ok := v.([]string)
			if ok && len(value) > 0 {
				query = query.Where(k+" IN ?", value)

			}
		case "status":
			value, ok := v.([]int)
			if ok && len(value) > 0 {
				query = query.Where("courier_service.status IN ?", value)

			}

		}
	}

	if len(sort) > 0 {
		if sort == "shipping_type_code" {
			sort = "shipping_type"
		}

		query = query.Order(sort)
	} else {
		query = query.Order("courier_service.updated_at DESC")
	}

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(couriers, &pagination, query, int64(len(couriers)))).
		Find(&couriersResponse).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return couriersResponse, &pagination, nil
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
