package repository

import (
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/helper/global"
	"strings"

	"gorm.io/gorm"
)

type ShippingCourierStatusRepository interface {
	FindByParams(limit int, page int, sort string, filters map[string]interface{}) ([]entity.ShippingCourierStatus, *base.Pagination, error)
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
		Joins("ShippingStatus").
		Joins("INNER JOIN channel ch ON ch.ID = \"ShippingStatus\".channel_id").
		Joins("Courier")

	for k, v := range filters {
		switch k {
		case "channel_name", "courier_name":
			value := v.([]string)
			if len(value) > 0 {
				query = query.Where(fmt.Sprint(k, " IN ?"), value)
			}
		case "status_code":
			value := v.([]string)
			if len(value) > 0 {
				query = query.Where("\"ShippingStatus\".status_code IN ?", value)
			}
		case "status_name":
			value, ok := v.([]string)
			if ok && len(value) > 0 {
				query = query.Where(global.AddLike(k, value))

			}
		case "status_courier":
			value := v.([]string)
			if len(value) > 0 {
				query = query.Where(global.AddLike("(shipping_courier_status.status_courier->'status')::text", value))
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
		sort = strings.ReplaceAll(sort, "status_title", "status_name")
		sort = strings.ReplaceAll(sort, "courier_status", "status_courier")
		query = query.Order(sort)
	} else {
		query = query.Order("shipping_courier_status.updated_at DESC")
	}

	pagination.Limit = limit
	pagination.Page = page

	err = query.Offset(pagination.GetOffset()).Limit((pagination.GetLimit())).Find(&result).Error
	if err != nil {
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
