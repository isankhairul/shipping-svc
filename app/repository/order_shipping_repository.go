package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"
	"go-klikdokter/pkg/util"
	"strings"

	"gorm.io/gorm"
)

type OrderShippingRepository interface {
	Create(input *entity.OrderShipping) (*entity.OrderShipping, error)
	Update(input *entity.OrderShipping) (*entity.OrderShipping, error)
	Upsert(input *entity.OrderShipping) (*entity.OrderShipping, error)
	FindByOrderNo(orderNo string) (*entity.OrderShipping, error)
	FindByUID(uid string) (*entity.OrderShipping, error)
	FindByParams(limit, page int, sort string, filter map[string]interface{}) ([]response.GetOrderShippingList, *base.Pagination, error)
	FindByUIDs(channelUID string, uid []string) ([]entity.OrderShipping, error)
}

type orderShippingRepository struct {
	base BaseRepository
}

func NewOrderShippingRepository(br BaseRepository) OrderShippingRepository {
	return &orderShippingRepository{br}
}

func (r *orderShippingRepository) Create(input *entity.OrderShipping) (*entity.OrderShipping, error) {
	db := r.base.GetDB()
	err := db.Create(input).Error

	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r *orderShippingRepository) Update(input *entity.OrderShipping) (*entity.OrderShipping, error) {
	db := r.base.GetDB()
	err := db.Updates(input).Error

	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r *orderShippingRepository) Upsert(input *entity.OrderShipping) (*entity.OrderShipping, error) {

	err := r.base.GetDB().Transaction(func(tx *gorm.DB) error {
		if input.ID == 0 {
			if err := tx.Create(input).Error; err != nil {
				return err
			}

			return nil
		}

		if err := tx.Updates(input).Error; err != nil {
			return err
		}

		return nil
	})

	return input, err
}

func (r *orderShippingRepository) FindByOrderNo(orderNo string) (*entity.OrderShipping, error) {
	var result entity.OrderShipping
	query := r.base.GetDB().
		Model(&entity.OrderShipping{}).
		Preload("Channel").
		Preload("Courier").
		Preload("CourierService").
		Preload("OrderShippingItem").
		Preload("OrderShippingHistory").
		Where(&entity.OrderShipping{OrderNo: orderNo})

	err := query.First(&result).Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *orderShippingRepository) FindByUID(uid string) (*entity.OrderShipping, error) {
	var result entity.OrderShipping
	query := r.base.GetDB().
		Preload("Channel").
		Preload("Courier").
		Preload("OrderShippingItem").
		Preload("CourierService").
		Preload("OrderShippingHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_shipping_history.created_at DESC")
		}).
		Preload("OrderShippingHistory.ShippingCourierStatus.ShippingStatus").
		Model(&entity.OrderShipping{}).
		Where(&entity.OrderShipping{BaseIDModel: base.BaseIDModel{UID: uid}})

	err := query.First(&result).Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *orderShippingRepository) FindByParams(limit, page int, sort string, filter map[string]interface{}) ([]response.GetOrderShippingList, *base.Pagination, error) {
	pagination := &base.Pagination{}

	var result []response.GetOrderShippingList

	query := r.base.GetDB().
		Model(&entity.OrderShipping{}).
		Select(
			"ch.channel_code AS channel_code",
			"ch.channel_name AS channel_name",
			"order_shipping.uid AS order_shipping_uid",
			"order_shipping.order_shipping_date AS order_shipping_date",
			"order_shipping.order_no AS order_no",
			"c.courier_name AS courier_name",
			"cs.shipping_name AS courier_services_name",
			"order_shipping.airwaybill AS airwaybill",
			"order_shipping.booking_id AS booking_id",
			"order_shipping.merchant_name AS merchant_name",
			"order_shipping.customer_name AS customer_name",
			"order_shipping.status AS shipping_status",
			"ss.status_name AS shipping_status_name",
		).
		Joins("INNER JOIN channel ch ON ch.id = order_shipping.channel_id").
		Joins("INNER JOIN courier c ON c.id = order_shipping.courier_id").
		Joins("INNER JOIN courier_service cs ON cs.id = order_shipping.courier_service_id").
		Joins("INNER JOIN shipping_status ss ON ss.status_code = order_shipping.status")

	for k, v := range filter {

		if !util.IsNilOrEmpty(v) {

			switch k {
			case "order_no":
				query = query.Where(like(k, v.([]string)))

			case "courier_name":
				query = query.Where(like("c.courier_name", v.([]string)))

			case "channel_name":
				query = query.Where(like("ch.channel_name", v.([]string)))

			case "channel_code":
				query = query.Where("ch.channel_code IN ?", v.([]string))

			case "shipping_status":
				query = query.Where("order_shipping.status IN ?", v.([]string))

			case "order_shipping_date_from":
				query = query.Where("CAST(order_shipping_date AS DATE) >= CAST(? AS DATE)", v)

			case "order_shipping_date_to":
				query = query.Where("CAST(order_shipping_date AS DATE) <= CAST(? AS DATE)", v)

			case "courier_services_name":
				query = query.Where(like("cs.shipping_name", v.([]string)))

			case "airwaybill":
				query = query.Where(like("order_shipping.airwaybill", v.([]string)))

			case "booking_id":
				query = query.Where(like("order_shipping.booking_id", v.([]string)))
			case "merchant_name":
				query = query.Where(like("order_shipping.merchant_name", v.([]string)))
			case "customer_name":
				query = query.Where(like("order_shipping.customer_name", v.([]string)))

			}

		}
	}

	sort = strings.ReplaceAll(sort, "courier_code", "c.code")
	sort = strings.ReplaceAll(sort, "shipping_status", "order_shipping.status")
	sort = strings.ReplaceAll(sort, "courier_services_name", "cs.shipping_name")
	sort = strings.ReplaceAll(sort, "order_shipping_uid", "order_shipping.uid")
	sort = strings.ReplaceAll(sort, "order_shipping_date", "order_shipping.order_shipping_date")

	sort = util.ReplaceEmptyString(sort, "order_shipping.updated_at desc")

	query = query.Order(sort)

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.base.Paginate(&entity.OrderShipping{}, pagination, query, int64(len(result)))).
		Find(&result).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return result, pagination, nil
}

func (r *orderShippingRepository) FindByUIDs(channelUID string, uid []string) ([]entity.OrderShipping, error) {
	var result []entity.OrderShipping
	query := r.base.GetDB().
		Preload("Channel").
		Preload("Courier").
		Preload("OrderShippingItem").
		Preload("CourierService").
		Model(&entity.OrderShipping{}).
		Where("order_shipping.uid IN ?", uid).
		Joins("INNER JOIN channel c ON c.id = order_shipping.channel_id AND c.uid = ?", channelUID)

	err := query.Find(&result).Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}
