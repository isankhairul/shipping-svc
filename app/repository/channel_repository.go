package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"math"

	"gorm.io/gorm"
)

type channelRepo struct {
	base BaseRepository
}

type ChannelRepository interface {
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Channel, *base.Pagination, error)
	FindByUid(uid *string) (*entity.Channel, error)
	CreateChannel(product *entity.Channel) (*entity.Channel, error)
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB
	Delete(uid string) error
	Update(uid string, input map[string]interface{}) error
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

func (r *channelRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Channel, *base.Pagination, error) {
	var channels []entity.Channel
	var pagination base.Pagination

	query := r.base.GetDB()

	if filter["channel_code"] != "" {
		query = query.Where("channel_code = ?", filter["channel_code"])
	}
	if filter["channel_name"] != "" {
		query = query.Where("channel_name = ?", filter["channel_name"])
	}
	if filter["status"] != 0 {
		query = query.Where("status = ?", filter["status"])
	}

	if len(sort) > 0 {
		query = query.Order(sort)
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
