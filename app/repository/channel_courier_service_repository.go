package repository

import (
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"strings"

	"gorm.io/gorm"
)

type ChannelCourierServiceRepositoryImpl struct {
	base BaseRepository
}

type ChannelCourierServiceRepository interface {
	GetChannelCourierService(channelCourierID, courierServiceID uint64) (*entity.ChannelCourierService, error)
	GetChannelCourierServiceByUID(uid string) (*entity.ChannelCourierService, error)
	//GetChannelCourierServicesByCourierServiceUIds(uids []*string) ([]*entity.ChannelCourierService, error)
	CreateChannelCourierService(input *entity.ChannelCourierService) (*entity.ChannelCourierService, error)
	UpdateChannelCourierService(uid string, updates map[string]interface{}) error
	DeleteChannelCourierServiceByID(id uint64) error
	DeleteChannelCourierServicesByChannelID(channelID uint64, courierID uint64) error
	FindByParams(limit, page int, sort string, filters map[string]interface{}) ([]entity.ChannelCourierService, *base.Pagination, error)
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
		Preload("ChannelCourier.Channel").
		Preload("ChannelCourier.Courier").
		Preload("CourierService").
		Joins("ChannelCourier").
		Joins("INNER JOIN channel ch ON ch.ID = \"ChannelCourier\".channel_id").
		Joins("INNER JOIN courier co ON co.ID = \"ChannelCourier\".courier_id").
		Joins("CourierService")

	for k, v := range filters {
		switch k {
		case "shipping_name", "shipping_code", "shipping_type", "channel_name", "courier_name":
			value := v.([]string)
			if len(value) > 0 {
				query = query.Where(fmt.Sprint(k, " IN ?"), value)
			}
		case "status":
			value := v.([]int)
			if len(value) > 0 {
				query = query.Where("channel_courier_service.status IN ? ", value)
			}
		}
	}

	var totalRecords int64
	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, nil, err
	}

	pagination.SetTotalRecords(totalRecords)

	if len(sort) > 0 {
		if strings.Contains(strings.ToLower(sort), "status") {
			sort = strings.ReplaceAll(sort, "status", "channel_courier_service.status")
		}
		query = query.Order(sort)
	} else {
		query = query.Order("channel_courier_service.updated_at DESC")
	}

	pagination.Limit = limit
	pagination.Page = page

	err = query.Offset(pagination.GetOffset()).Limit((pagination.GetLimit())).Find(&result).Error
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
