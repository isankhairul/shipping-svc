package repository

import (
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"gorm.io/gorm"
)

type ChannelCourierRepositoryImpl struct {
	base BaseRepository
}

type ChannelCourierRepository interface {
	CreateChannelCourier(item *entity.ChannelCourier) (*entity.ChannelCourier, error)
	GetChannelCourierByIds(channelID uint64, courierID uint64) (*entity.ChannelCourier, error)
	GetChannelCourierByUID(uid string) (*entity.ChannelCourier, error)
	FindByPagination(limit int, page int, sort string, filters map[string]interface{}) ([]*entity.ChannelCourier, *base.Pagination, error)
	DeleteChannelCourierByID(id uint64) error
	UpdateChannelCourier(uid string, updates map[string]interface{}) error
	FindCourierByUID(uid string) (*entity.Courier, error)
	FindChannelByUID(uid string) (*entity.Channel, error)
	IsHasChannelCourierService(channelCourierID uint64) bool
}

func NewChannelCourierRepository(br BaseRepository) ChannelCourierRepository {
	return &ChannelCourierRepositoryImpl{br}
}

func (s *ChannelCourierRepositoryImpl) UpdateChannelCourier(uid string, updates map[string]interface{}) error {
	result := s.base.GetDB().Model(&entity.ChannelCourier{}).
		Where(&entity.ChannelCourier{BaseIDModel: base.BaseIDModel{UID: uid}}).
		Updates(updates)
	if result != nil && result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}

func (s *ChannelCourierRepositoryImpl) FindCourierByUID(courierUID string) (*entity.Courier, error) {
	db := s.base.GetDB()
	var courier *entity.Courier
	notFoundCourier := db.Where(&entity.Courier{BaseIDModel: base.BaseIDModel{UID: courierUID}}).First(&courier).Error
	if notFoundCourier != nil {
		return nil, notFoundCourier
	}
	return courier, nil
}

func (s *ChannelCourierRepositoryImpl) FindChannelByUID(channelUID string) (*entity.Channel, error) {
	db := s.base.GetDB()
	var channel *entity.Channel
	notFound := db.Where(&entity.Channel{BaseIDModel: base.BaseIDModel{UID: channelUID}}).First(&channel).Error
	if notFound != nil {
		return nil, notFound
	}
	return channel, nil
}

func (r *ChannelCourierRepositoryImpl) CreateChannelCourier(channelCourier *entity.ChannelCourier) (*entity.ChannelCourier, error) {
	err := r.base.GetDB().Create(channelCourier).Error
	if err != nil {
		return nil, err
	}
	return channelCourier, nil
}

func (r *ChannelCourierRepositoryImpl) GetChannelCourierByIds(channelID uint64, courierID uint64) (*entity.ChannelCourier, error) {
	var cur entity.ChannelCourier
	db := r.base.GetDB()
	err := db.Preload("Courier").Preload("Channel").
		Where(&entity.ChannelCourier{ChannelID: channelID, CourierID: courierID}).First(&cur).Error
	if err != nil {
		return nil, err
	}
	return &cur, nil
}

func (r *ChannelCourierRepositoryImpl) GetChannelCourierByUID(uid string) (*entity.ChannelCourier, error) {
	var cur *entity.ChannelCourier
	db := r.base.GetDB()

	err := db.Preload("Courier").Preload("Channel").
		Where(&entity.ChannelCourier{BaseIDModel: base.BaseIDModel{UID: uid}}).First(&cur).Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}

	return cur, nil
}

func (r *ChannelCourierRepositoryImpl) DeleteChannelCourierByID(id uint64) error {
	cur := &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{ID: id}}
	db := r.base.GetDB()
	err := db.Model(&entity.ChannelCourier{}).Where(&entity.ChannelCourier{BaseIDModel: base.BaseIDModel{ID: id}}).Delete(&cur).Error
	return err
}

func (r *ChannelCourierRepositoryImpl) IsHasChannelCourierService(channelCourierID uint64) bool {
	var count int64
	r.base.GetDB().Model(&entity.ChannelCourierService{}).
		Where(&entity.ChannelCourierService{ChannelCourierID: channelCourierID}).Count(&count)

	return count > 0
}

func (r *ChannelCourierRepositoryImpl) FindByPagination(limit int, page int, sort string, filters map[string]interface{}) ([]*entity.ChannelCourier, *base.Pagination, error) {
	var items []*entity.ChannelCourier
	var pagination base.Pagination
	db := r.base.GetDB()
	query := db.Model(&entity.ChannelCourier{}).
		Preload("Channel").Preload("Courier").Joins("Courier").Joins("Channel")

	for k, v := range filters {
		switch k {
		case "channel_code", "channel_name", "courier_name":
			value := v.([]string)
			if len(value) > 0 {
				query = query.Where(fmt.Sprint(k, " IN ?"), value)
			}
		case "status":
			value := v.([]int)
			if len(value) > 0 {
				query = query.Where("channel_courier.status IN ? ", value)
			}
		}
	}

	// sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Model(&entity.ChannelCourier{}).
	// 		Joins("Courier").
	// 		Joins("Channel").
	// 		Find(&[]entity.ChannelCourier{})
	// })

	// println(sql)

	var totalRecords int64
	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, nil, err
	}

	pagination.SetTotalRecords(totalRecords)

	if len(sort) > 0 {
		query = query.Order(sort)
	} else {
		query = query.Order("updated_at DESC")
	}

	pagination.Limit = limit
	pagination.Page = page

	err = query.Offset(pagination.GetOffset()).Limit((pagination.GetLimit())).Find(&items).Error
	if err != nil {
		return nil, nil, err
	}

	return items, &pagination, nil
}
