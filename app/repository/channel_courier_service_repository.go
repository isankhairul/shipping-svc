package repository

import (
	"go-klikdokter/app/model/entity"
)

type ChannelCourierServiceRepositoryImpl struct {
	base BaseRepository
}

type ChannelCourierServiceRepository interface {
	GetChannelCourierServicesByCourierServiceUIds(uids []*string) ([]*entity.ChannelCourierService, error)
	CreateChannelCourierService(courier *entity.Courier, channel *entity.Channel, cs *entity.CourierService, priceInternal float64, status int) (*entity.ChannelCourierService, error)
	DeleteChannelCourierServiceByID(id uint64) error
	DeleteChannelCourierServicesByChannelID(channelID uint64, courierID uint64) error
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

func (r *ChannelCourierServiceRepositoryImpl) CreateChannelCourierService(courier *entity.Courier, channel *entity.Channel, cs *entity.CourierService,
	priceInternal float64, status int) (*entity.ChannelCourierService, error) {
	db := r.base.GetDB()
	var cur *entity.ChannelCourierService

	_ = db.Model(&entity.ChannelCourierService{}).
		Where(&entity.ChannelCourierService{
			CourierID:        courier.ID,
			ChannelID:        channel.ID,
			CourierServiceID: cs.ID}).First(&cur).Error
	if cur.ID > 0 {
		changed := cur.PriceInternal != priceInternal || cur.Status != status
		if cur.PriceInternal != priceInternal {
			cur.PriceInternal = priceInternal
		}
		if cur.Status != status {
			cur.Status = status
		}
		if changed {
			err := db.Save(&cur).Error
			if err != nil {
				return nil, err
			}
		}
		return cur, nil
	}

	cur = &entity.ChannelCourierService{
		CourierID:        courier.ID,
		ChannelID:        channel.ID,
		CourierServiceID: cs.ID,
		PriceInternal:    priceInternal,
		Status:           status,
	}
	err := db.Model(&entity.ChannelCourierService{}).Create(&cur).Error
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (r *ChannelCourierServiceRepositoryImpl) DeleteChannelCourierServiceByID(uid uint64) error {
	db := r.base.GetDB()
	var ret entity.ChannelCourierService
	err := db.Model(&entity.ChannelCourierService{}).Where("id=?", uid).Delete(&ret).Error
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
