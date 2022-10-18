package repository

import (
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"
	"go-klikdokter/pkg/util"
	"strings"

	"gorm.io/gorm"
)

type ChannelCourierServiceRepositoryImpl struct {
	base BaseRepository
}

const (
	channelCourierServiceStatus = "channel_courier_service.status"
)

type ChannelCourierServiceRepository interface {
	GetChannelCourierService(channelCourierID, courierServiceID uint64) (*entity.ChannelCourierService, error)
	GetChannelCourierServiceByUID(uid string) (*entity.ChannelCourierService, error)
	//GetChannelCourierServicesByCourierServiceUIds(uids []*string) ([]*entity.ChannelCourierService, error)
	CreateChannelCourierService(input *entity.ChannelCourierService) (*entity.ChannelCourierService, error)
	UpdateChannelCourierService(uid string, updates map[string]interface{}) error
	DeleteChannelCourierServiceByID(id uint64) error
	DeleteChannelCourierServicesByChannelID(channelID uint64, courierID uint64) error
	FindByParams(limit, page int, sort string, filters map[string]interface{}) ([]entity.ChannelCourierService, *base.Pagination, error)
	GetChannelCourierListByChannelUID(channelUID string, limit int, page int, sort, dir string, filter map[string]interface{}) ([]response.CourierServiceByChannelResponse, *base.Pagination, error)
}

func NewChannelCourierServiceRepository(br BaseRepository) ChannelCourierServiceRepository {
	return &ChannelCourierServiceRepositoryImpl{br}
}

func (r *ChannelCourierServiceRepositoryImpl) GetChannelCourierServicesByCourierServiceUIds(uids []*string) ([]*entity.ChannelCourierService, error) {
	items := []*entity.ChannelCourierService{}
	err := r.base.GetDB().Model(&entity.ChannelCourierService{}).Where("courier_service_uid IN (?)", uids).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ChannelCourierServiceRepositoryImpl) CreateChannelCourierService(input *entity.ChannelCourierService) (*entity.ChannelCourierService, error) {
	db := r.base.GetDB()

	err := db.Model(&entity.ChannelCourierService{}).
		Create(&input).
		Error

	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r *ChannelCourierServiceRepositoryImpl) GetChannelCourierServiceByUID(uid string) (*entity.ChannelCourierService, error) {
	var result entity.ChannelCourierService

	err := r.base.GetDB().
		Model(&entity.ChannelCourierService{}).
		Preload("ChannelCourier.Channel").
		Preload("ChannelCourier.Courier").
		Preload("CourierService").
		Where("uid=?", uid).
		First(&result).
		Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *ChannelCourierServiceRepositoryImpl) FindByParams(limit, page int, sort string, filters map[string]interface{}) ([]entity.ChannelCourierService, *base.Pagination, error) {
	var pagination base.Pagination
	var result []entity.ChannelCourierService
	db := r.base.GetDB()

	query := db.Model(&entity.ChannelCourierService{}).
		Select("channel_courier_service.*", "st.title AS shipping_type_name").
		Preload("ChannelCourier.Channel").
		Preload("ChannelCourier.Courier").
		Preload("CourierService").
		Joins("ChannelCourier").
		Joins("INNER JOIN channel ch ON ch.ID = \"ChannelCourier\".channel_id").
		Joins("INNER JOIN courier co ON co.ID = \"ChannelCourier\".courier_id").
		Joins("INNER JOIN courier_service cs ON cs.id = channel_courier_service.courier_service_id").
		Joins("INNER JOIN shippment_predefined st ON st.code = cs.shipping_type AND st.type = 'shipping_type'")

	for k, v := range filters {
		k = strings.ReplaceAll(k, "courier_uid", "co.uid")
		k = strings.ReplaceAll(k, "status", channelCourierServiceStatus)
		k = strings.ReplaceAll(k, "shipping_type_name", "st.title")

		if util.IsSliceAndNotEmpty(v) {
			query = query.Where(fmt.Sprint(k, " IN ?"), v)
		}
	}

	sort = strings.ReplaceAll(strings.ToLower(sort), "status", channelCourierServiceStatus)

	if len(sort) == 0 {
		sort = "channel_courier_service.updated_at DESC"
	}

	query = query.Order(sort)

	pagination.Limit = limit
	pagination.Page = page

	err := query.Scopes(r.Paginate(result, &pagination, query, int64(len(result)))).
		Find(&result).
		Error

	if err != nil {
		return nil, nil, err
	}

	return result, &pagination, nil
}

func (s *ChannelCourierServiceRepositoryImpl) UpdateChannelCourierService(uid string, updates map[string]interface{}) error {
	err := s.base.GetDB().Model(&entity.ChannelCourierService{}).
		Where(&entity.ChannelCourierService{BaseIDModel: base.BaseIDModel{UID: uid}}).
		Updates(updates).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *ChannelCourierServiceRepositoryImpl) DeleteChannelCourierServiceByID(id uint64) error {
	db := r.base.GetDB()
	var ret entity.ChannelCourierService
	err := db.Model(&entity.ChannelCourierService{}).Where("id=?", id).Delete(&ret).Error
	return err
}

func (r *ChannelCourierServiceRepositoryImpl) DeleteChannelCourierServicesByChannelID(channelID uint64, courierID uint64) error {
	db := r.base.GetDB()
	var ret entity.ChannelCourierService
	err := db.Model(&entity.ChannelCourierService{}).
		Where("channel_id=? AND courier_id=?", channelID, courierID).
		Delete(&ret).Error
	return err
}

func (r *ChannelCourierServiceRepositoryImpl) GetChannelCourierService(channelCourierID, courierServiceID uint64) (*entity.ChannelCourierService, error) {
	var result entity.ChannelCourierService

	err := r.base.GetDB().
		Model(&entity.ChannelCourierService{}).
		Preload("ChannelCourier").
		Preload("CourierService").
		Where("channel_courier_id=? AND courier_service_id=?", channelCourierID, courierServiceID).
		First(&result).
		Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *ChannelCourierServiceRepositoryImpl) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)
	pagination.SetTotalRecords(totalRecords)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *ChannelCourierServiceRepositoryImpl) GetChannelCourierListByChannelUID(channelUID string, limit int, page int, sort, dir string, filter map[string]interface{}) ([]response.CourierServiceByChannelResponse, *base.Pagination, error) {
	db := r.base.GetDB()
	var courierService []response.CourierServiceByChannelResponse
	var pagination base.Pagination
	query := db.Model(&entity.ChannelCourierService{}).
		Select(
			"cs.uid AS courier_service_uid",
			"cs.shipping_code AS shipping_code",
			"cs.shipping_name AS shipping_name",
			"cs.shipping_description AS shipping_description",
			"cs.image_path AS image_logo",
			"cs.etd_min AS etd_min",
			"cs.etd_max AS etd_max",
			"cs.shipping_type AS shipping_type_code",
			"st.title AS shipping_type_name",
			"c.uid AS courier_uid",
			"c.code AS courier_code",
			"c.courier_name AS courier_name",
			"c.courier_type AS courier_type_code",
			"ct.title AS courier_type_name",
			"c.image_path AS courier_image",
		).
		Joins("INNER JOIN channel_courier cc ON cc.id = channel_courier_service.channel_courier_id").
		Joins("INNER JOIN courier_service cs ON cs.id = channel_courier_service.courier_service_id").
		Joins("INNER JOIN channel ch ON ch.id = cc.channel_id").
		Joins("INNER JOIN courier c ON cc.courier_id = c.id").
		Joins("LEFT JOIN shippment_predefined ct ON ct.code = c.courier_type AND ct.type = 'courier_type'").
		Joins("LEFT JOIN shippment_predefined st ON st.code = cs.shipping_type AND st.type = 'shipping_type'").
		Where("ch.uid = ?", channelUID)

	for k, v := range filter {

		if util.IsSliceAndNotEmpty(v) {

			k = strings.ReplaceAll(k, "courier_type_code", "ct.code")
			k = strings.ReplaceAll(k, "courier_code", "c.code")
			k = strings.ReplaceAll(k, "courier_name", "c.courier_name")
			k = strings.ReplaceAll(k, "shipping_type_code", "st.code")
			k = strings.ReplaceAll(k, "shipping_name", "cs.shipping_name")
			k = strings.ReplaceAll(k, "status", channelCourierServiceStatus)

			query = query.Where(fmt.Sprint(k, " IN ?"), v)
		}
	}

	if len(sort) == 0 {
		sort = "cs.id"
	}

	if strings.EqualFold(dir, "desc") {
		sort += " desc"
	}

	query = query.Order(sort)
	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate([]entity.ChannelCourierService{}, &pagination, query, int64(len([]entity.ChannelCourierService{})))).
		Find(&courierService).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	for i := range courierService {
		courierService[i].Courier = response.CourierByChannelResponse{
			CourierUID:      courierService[i].CourierUID,
			CourierCode:     courierService[i].CourierCode,
			CourierName:     courierService[i].CourierName,
			CourierTypeCode: courierService[i].CourierTypeCode,
			CourierTypeName: courierService[i].CourierTypeName,
			ImageLogo:       courierService[i].CourierImage,
		}
	}

	return courierService, &pagination, nil
}
