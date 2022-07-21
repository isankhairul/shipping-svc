package repository

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
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
		return nil, err
	}
	var items []*entity.ChannelCourierService
	err = db.Model(&entity.ChannelCourierService{}).Preload("CourierService").
		Where(&entity.ChannelCourierService{ChannelID: cur.ChannelID, CourierID: cur.CourierID}).Find(&items).Error
	if items != nil {
		cur.ChannelCourierServices = items
	}
	return cur, nil
}

func (r *ChannelCourierRepositoryImpl) DeleteChannelCourierByID(id uint64) error {
	cur := &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{ID: id}}
	db := r.base.GetDB()
	err := db.Model(&entity.ChannelCourier{}).Where(&entity.ChannelCourier{BaseIDModel: base.BaseIDModel{ID: id}}).Delete(&cur).Error
	return err
}

func (r *ChannelCourierRepositoryImpl) FindByPagination(limit int, page int, sort string, filters map[string]interface{}) ([]*entity.ChannelCourier, *base.Pagination, error) {
	var items []*entity.ChannelCourier
	var pagination base.Pagination
	db := r.base.GetDB()
	query := db.Model(&entity.ChannelCourier{}).
		Preload("Channel").Preload("Courier").Joins("Courier").Joins("Channel")

	if filters["channel_name"] != "" {
		query = query.Where(&entity.ChannelCourier{Channel: &entity.Channel{ChannelName: filters["channel_name"].(string)}})
	}

	if filters["channel_code"] != "" {
		query = query.Where(&entity.ChannelCourier{Channel: &entity.Channel{ChannelCode: filters["channel_code"].(string)}})
	}

	if filters["courier_name"] != "" {
		query = query.Where(&entity.ChannelCourier{Courier: &entity.Courier{CourierName: filters["courier_name"].(string)}})
	}

	if filters["status"].(*int) != nil {
		query = query.Where(&entity.ChannelCourier{Status: *filters["status"].(*int)})
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
