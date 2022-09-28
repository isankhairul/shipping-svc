package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"
	"go-klikdokter/pkg/util"
	"math"
	"strings"

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
	FindCourierServiceByChannelAndUIDs(channel_uid string, uids []string, containPrescription bool, shippingType string) ([]entity.ChannelCourierServiceForShippingRate, error)
	FindCourierService(channelUID, courierServiceUID string) (*entity.CourierService, error)
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

		if util.IsSliceAndNotEmpty(v) {

			switch k {
			case "shipping_code", "shipping_name":
				query = query.Where(like(k, v.([]string)))

			case "courier_uid", "courier_type", "shipping_type":
				query = query.Where(k+" IN ?", v.([]string))

			case "status":
				query = query.Where("courier_service.status IN ?", v)

			}
		}
	}

	if strings.Contains(sort, "shipping_type_code") {
		m := map[string]string{"shipping_type_code": "shipping_type", "shipping_type_code asc": "shipping_type asc", "shipping_type_code desc": "shipping_type desc"}
		sort = m[strings.ToLower(sort)]
	}

	if len(sort) == 0 {
		sort = "courier_service.updated_at DESC"
	}

	query = query.Order(sort)

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

func (r *courierServiceRepo) FindCourierServiceByChannelAndUIDs(channel_uid string, uids []string, containPrescription bool, shippingType string) ([]entity.ChannelCourierServiceForShippingRate, error) {

	db := r.base.GetDB()
	var courierService []entity.ChannelCourierServiceForShippingRate

	query := db.Model(&entity.ChannelCourierService{}).
		Select(
			"cs.uid AS courier_service_uid",
			"cs.shipping_code AS shipping_code",
			"cs.shipping_name AS shipping_name",
			"cs.shipping_description AS shipping_description",
			"cs.image_path AS logo",
			"cs.etd_min AS etd_min",
			"cs.etd_max AS etd_max",
			"cs.shipping_type AS shipping_type_code",
			"st.title AS shipping_type_name",
			"st.note AS shipping_type_description",
			"c.id AS courier_id",
			"c.uid AS courier_uid",
			"c.code AS courier_code",
			"c.courier_name AS courier_name",
			"c.courier_type AS courier_type_code",
			"ct.title AS courier_type_name",
			"cs.insurance_fee AS insurance_fee",
			"cs.insurance AS use_insurance",
			"channel_courier_service.price_internal AS price",
			"cs.max_weight AS max_weight",
			"c.status AS courier_status",
			"cs.status AS courier_service_status",
			"cc.status AS channel_courier_status",
			"channel_courier_service.status AS channel_courier_service_status",
			"c.hide_purpose AS hide_purpose",
			"cs.prescription_allowed AS prescription_allowed",
		).
		Joins("INNER JOIN channel_courier cc ON cc.id = channel_courier_service.channel_courier_id").
		Joins("INNER JOIN courier_service cs ON cs.id = channel_courier_service.courier_service_id").
		Joins("INNER JOIN channel ch ON ch.id = cc.channel_id").
		Joins("INNER JOIN courier c ON cc.courier_id = c.id").
		Joins("INNER JOIN shippment_predefined ct ON ct.code = c.courier_type AND ct.type = 'courier_type'").
		Joins("INNER JOIN shippment_predefined st ON st.code = cs.shipping_type AND st.type = 'shipping_type'").
		Where("ch.uid = ?", channel_uid).
		Where("cs.uid IN ?", uids)
		/*
				Where("channel_courier_service.status = 1").
				Where("cc.status = 1").
				Where("c.status = 1").
				Where("cs.status = 1").
				Where("c.hide_purpose = 0")

			if containPrescription {
				query = query.Where("cs.prescription_allowed = 1")
			}
		*/

	if shippingType != "" {
		query = query.Where("cs.shipping_type = ?", shippingType)
	}

	err := query.Find(&courierService).Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}

	return courierService, nil
}

func (r *courierServiceRepo) FindCourierService(channelUID, courierServiceUID string) (*entity.CourierService, error) {
	db := r.base.GetDB()
	var courierService *entity.CourierService

	query := db.Model(&entity.CourierService{}).
		Preload("Courier").
		Joins("INNER JOIN channel_courier_service ccs ON ccs.courier_service_id = courier_service.id").
		Joins("INNER JOIN channel_courier cc ON cc.id = ccs.channel_courier_id").
		Joins("INNER JOIN channel ch ON ch.id = cc.channel_id").
		Joins("INNER JOIN courier c ON cc.courier_id = courier_service.courier_id").
		Where("courier_service.uid = ?", courierServiceUID).
		Where("ch.uid = ?", channelUID)

	err := query.First(&courierService).Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}

	return courierService, nil
}
