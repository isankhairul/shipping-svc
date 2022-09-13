package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"gorm.io/gorm"
)

type OrderShippingRepository interface {
	Create(input *entity.OrderShipping) (*entity.OrderShipping, error)
	Update(input *entity.OrderShipping) (*entity.OrderShipping, error)
	Upsert(input *entity.OrderShipping) (*entity.OrderShipping, error)
	FindByOrderNo(orderNo string) (*entity.OrderShipping, error)
	FindByUID(uid string) (*entity.OrderShipping, error)
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
