package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"
	"gorm.io/gorm"
)

type courierRepo struct {
	base BaseRepository
}

type CourierRepository interface {
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]response.CourierListResponse, *base.Pagination, error)
	FindByUid(uid *string) (*entity.Courier, error)
	FindByCode(code string) (*entity.Courier, error)
	CreateCourier(courier *entity.Courier) (*entity.Courier, error)
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB
	Delete(uid string) error
	Update(uid string, input map[string]interface{}) error
	IsCourierHasChild(courierID uint64) *entity.CourierHasChildFlag
}

func NewCourierRepository(br BaseRepository) CourierRepository {
	return &courierRepo{br}
}

func (r *courierRepo) FindByUid(uid *string) (*entity.Courier, error) {
	var courier entity.Courier
	err := r.base.GetDB().Preload("CourierServices").
		Where(&entity.Courier{BaseIDModel: base.BaseIDModel{UID: *uid}}).
		First(&courier).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}

	return &courier, nil
}

func (r *courierRepo) FindByCode(code string) (*entity.Courier, error) {
	var courier entity.Courier
	err := r.base.GetDB().Where("code=?", code).First(&courier).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
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
	pagination.SetTotalRecords(totalRecords)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *courierRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]response.CourierListResponse, *base.Pagination, error) {
	var couriers []*entity.Courier
	var pagination base.Pagination
	var respCourier []response.CourierListResponse

	query := r.base.GetDB().
		Model(&entity.Courier{}).
		Select("courier.*, sp.title as courier_type_name").
		Joins("left join shippment_predefined as sp on courier.courier_type = sp.code").
		Where("sp.type = 'courier_type'")

	for k, v := range filter {
		switch k {
		case "courier_name":
			value, ok := v.([]string)
			if ok && len(value) > 0 {
				query = query.Where(like(k, value))

			}
		case "code":
			value, ok := v.([]string)
			if ok && len(value) > 0 {
				query = query.Where(like("courier.code", value))

			}
		case "courier_type":
			value, ok := v.([]string)
			if ok && len(value) > 0 {
				query = query.Where("courier_type IN ?", value)

			}
		case "status":
			value, ok := v.([]int)
			if ok && len(value) > 0 {
				query = query.Where("courier.status IN ?", value)

			}

		}
	}

	if len(sort) > 0 {
		m := map[string]string{"courier_code": "courier.code", "courier_type_name": "courier.courier_type_name"}
		sortValue := m[sort]
		query = query.Order(sortValue)
	} else {
		query = query.Order("updated_at DESC")
	}

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(couriers, &pagination, query, int64(len(couriers)))).
		Find(&respCourier).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return respCourier, &pagination, nil
}

func (r *courierRepo) Delete(uid string) error {

	var courier entity.Courier
	err := r.base.GetDB().
		Where("uid = ?", uid).First(&courier).Error
	if err != nil {
		return err
	}

	err = r.base.GetDB().
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

func (r *courierRepo) IsCourierHasChild(courierID uint64) *entity.CourierHasChildFlag {
	db := r.base.GetDB()
	var courierService int64
	var courierCoverage int64
	var channelCourier int64
	var shippingStatus int64

	db.Model(&entity.CourierService{}).Where(&entity.CourierService{CourierID: courierID}).Count(&courierService)
	db.Model(&entity.CourierCoverageCode{}).Where(&entity.CourierCoverageCode{CourierID: courierID}).Count(&courierCoverage)
	db.Model(&entity.ChannelCourier{}).Where(&entity.ChannelCourier{CourierID: courierID}).Count(&channelCourier)
	db.Model(&entity.ShippingCourierStatus{}).Where(&entity.ShippingCourierStatus{CourierID: courierID}).Count(&shippingStatus)

	return &entity.CourierHasChildFlag{
		CourierService:        courierService > 0,
		CourierCoverageCode:   courierCoverage > 0,
		ChannelCourier:        channelCourier > 0,
		ShippingCourierStatus: shippingStatus > 0,
	}
}
