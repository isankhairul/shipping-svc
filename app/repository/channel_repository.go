package repository

import (
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"math"
	"strings"

	"gorm.io/gorm"
)

type channelRepo struct {
	base BaseRepository
}

type ChannelRepository interface {
	FindAll(limit int, page int, sort string) ([]entity.Channel, *base.Pagination, error)
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Channel, *base.Pagination, error)
	CheckExistsByUIdChannelCode(uid, channelCode string) (bool, error)
	FindByUid(uid *string) (*entity.Channel, error)
	CreateChannel(channel *entity.Channel) (*entity.Channel, error)
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB
	Delete(uid string) error
	Update(uid string, input map[string]interface{}) error
	IsChannelHasChild(channelID uint64) *entity.ChannelHasChildFlag
}

func NewChannelRepository(br BaseRepository) ChannelRepository {
	return &channelRepo{br}
}

func (r *channelRepo) FindByUid(uid *string) (*entity.Channel, error) {
	var channel entity.Channel
	err := r.base.GetDB().
		Where("uid=?", uid).
		First(&channel).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}

	return &channel, nil
}

func (r *channelRepo) CreateChannel(Channel *entity.Channel) (*entity.Channel, error) {
	err := r.base.GetDB().
		Create(Channel).Error
	if err != nil {
		return nil, err
	}

	return Channel, nil
}

func (r *channelRepo) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(currRecord)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *channelRepo) FindAll(limit int, page int, sort string) ([]entity.Channel, *base.Pagination, error) {
	var channels []entity.Channel
	var pagination base.Pagination

	query := r.base.GetDB()
	if len(sort) > 0 {
		query = query.Order(sort)
	} else {
		query = query.Order("updated_at DESC")
	}
	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(channels, &pagination, query, int64(len(channels)))).
		Find(&channels).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return channels, &pagination, nil
}

func (r *channelRepo) CheckExistsByUIdChannelCode(uid, channelCode string) (bool, error) {
	var exists bool
	err := r.base.GetDB().
		Model(&entity.Channel{}).
		Select("count(*) > 0").
		Where("uid != ? AND channel_code = ?", uid, channelCode).
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

func (r *channelRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Channel, *base.Pagination, error) {
	var channels []entity.Channel
	var pagination base.Pagination

	query := r.base.GetDB()

	for k, v := range filter {
		switch k {
		case "channel_code", "channel_name":
			value, ok := v.([]string)
			if ok && len(value) > 0 {
				query = query.Where(like(k, value))

			}
		case "status":
			value, ok := v.([]int)
			if ok && len(value) > 0 {
				query = query.Where("status IN ?", value)

			}

		}
	}

	if len(sort) > 0 {
		query = query.Order(sort)
	} else {
		query = query.Order("updated_at DESC")
	}

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(channels, &pagination, query, int64(len(channels)))).
		Find(&channels).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return channels, &pagination, nil
}

func (r *channelRepo) Delete(uid string) error {
	var channel entity.Channel
	err := r.base.GetDB().
		Where("uid = ?", uid).
		Delete(&channel).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r *channelRepo) Update(uid string, input map[string]interface{}) error {
	err := r.base.GetDB().Model(&entity.Channel{}).
		Where("uid=?", uid).
		Updates(input).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *channelRepo) IsChannelHasChild(channelID uint64) *entity.ChannelHasChildFlag {
	db := r.base.GetDB()
	var channelCourier int64
	var shippingStatus int64

	db.Model(&entity.ChannelCourier{}).Where(&entity.ChannelCourier{ChannelID: channelID}).Count(&channelCourier)
	db.Model(&entity.ShippingStatus{}).Where(&entity.ShippingStatus{ChannelID: channelID}).Count(&shippingStatus)

	return &entity.ChannelHasChildFlag{
		ChannelCourier: channelCourier > 0,
		ShippingStatus: shippingStatus > 0,
	}
}

func like(column string, value []string) string {
	var condition string
	for _, v := range value {
		condition += fmt.Sprintf(" %s ILIKE '%%%s%%' OR", column, v)
	}
	return strings.TrimRight(condition, " OR")
}
