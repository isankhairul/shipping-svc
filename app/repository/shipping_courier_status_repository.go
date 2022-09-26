package repository

import (
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/helper/global"
	"go-klikdokter/pkg/util"
	"strings"

	"gorm.io/gorm"
)

type ShippingCourierStatusRepository interface {
	FindByParams(limit int, page int, sort string, filters map[string]interface{}) ([]entity.ShippingCourierStatus, *base.Pagination, error)
	FindByCode(channelID, courierID uint64, statusCode string) (*entity.ShippingCourierStatus, error)
	FindByCourierStatus(courierID uint64, statusCode string) (*entity.ShippingCourierStatus, error)
}

type shippingCourierStatusRepositoryImpl struct {
	base BaseRepository
}

func NewShippingCourierStatusRepository(br BaseRepository) ShippingCourierStatusRepository {
	return &shippingCourierStatusRepositoryImpl{br}
}

func (r *shippingCourierStatusRepositoryImpl) FindByParams(limit int, page int, sort string, filters map[string]interface{}) ([]entity.ShippingCourierStatus, *base.Pagination, error) {
	var pagination base.Pagination
	var result []entity.ShippingCourierStatus
	db := r.base.GetDB()

	query := db.Model(&entity.ShippingCourierStatus{}).
		Preload("ShippingStatus.Channel").
		Preload("Courier").
		Preload("ShippingStatus").
		Joins("INNER JOIN shipping_status ss ON ss.id = shipping_courier_status.shipping_status_id").
		Joins("INNER JOIN channel ch ON ch.ID = ss.channel_id").
		Joins("INNER JOIN courier co ON co.id = shipping_courier_status.courier_id")

	for k, v := range filters {

		if util.IsSliceAndNotEmpty(v) {

			switch k {
			case "channel_name", "courier_name", "channel_code":
				query = query.Where(fmt.Sprint(k, " IN ?"), v.([]string))

			case "status_code":
				query = query.Where("ss.status_code IN ?", v.([]string))

			case "status_name":
				query = query.Where(global.AddLike(k, v.([]string)))

			case "status_courier":
				query = query.Where(global.AddLike("(shipping_courier_status.status_courier->'status')::text", v.([]string)))

			}
		}
	}

	sort = strings.ReplaceAll(sort, "status_title", "status_name")
	sort = strings.ReplaceAll(sort, "courier_status", "status_courier")

	if len(sort) == 0 {
		sort = "shipping_courier_status.updated_at DESC"

	}

	query = query.Order(sort)

	pagination.Limit = limit
	pagination.Page = page

	err := query.Scopes(r.Paginate(result, &pagination, query, int64(len(result)))).
		Find(&result).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return result, &pagination, nil
}

func (r *shippingCourierStatusRepositoryImpl) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)
	pagination.SetTotalRecords(totalRecords)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *shippingCourierStatusRepositoryImpl) FindByCode(channelID, courierID uint64, statusCode string) (*entity.ShippingCourierStatus, error) {
	result := &entity.ShippingCourierStatus{}
	query := r.base.GetDB().
		Preload("ShippingStatus").
		Joins("INNER JOIN shipping_status ss ON ss.ID = shipping_courier_status.shipping_status_id AND ss.channel_id = ?", channelID).
		Where(&entity.ShippingCourierStatus{StatusCode: statusCode}).
		Where(&entity.ShippingCourierStatus{CourierID: courierID})

	err := query.First(result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return result, nil
}

func (r *shippingCourierStatusRepositoryImpl) FindByCourierStatus(courierID uint64, statusCode string) (*entity.ShippingCourierStatus, error) {
	result := &entity.ShippingCourierStatus{}
	query := r.base.GetDB().
		Where(&entity.ShippingCourierStatus{CourierID: courierID}).
		Where(fmt.Sprintf("(shipping_courier_status.status_courier->'status')::text ilike '%%\"%s\"%%'", statusCode))

	err := query.First(result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return result, nil
}
