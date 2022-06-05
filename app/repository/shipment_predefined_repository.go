package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"strings"

	"gorm.io/gorm"
)

type ShipmentPredefinedRepository interface {
	GetAll(limit int, page int, sort string, filter map[string]interface{}) ([]*entity.ShippmentPredefined, *base.Pagination, error)
	UpdateShipmentPredefined(dto entity.ShippmentPredefined) (*entity.ShippmentPredefined, error)
	GetShipmentPredefinedByUid(uid string) (*entity.ShippmentPredefined, error)
}

type ShipmentPredefinedRepositoryImpl struct {
	base BaseRepository
}

func NewShipmentPredefinedRepository(br BaseRepository) ShipmentPredefinedRepository {
	return &ShipmentPredefinedRepositoryImpl{br}
}

// GetAll implements ShipmentPredefinedRepository
func (r *ShipmentPredefinedRepositoryImpl) GetAll(limit int, page int, sort string, filter map[string]interface{}) ([]*entity.ShippmentPredefined, *base.Pagination, error) {
	var items []*entity.ShippmentPredefined
	var pagination base.Pagination

	query := r.base.GetDB().Model(&entity.ShippmentPredefined{})

	if filter["type"] != "" {
		query = query.Where("LOWER(type) LIKE ?", "%"+strings.ToLower(filter["type"].(string))+"%")
	}
	if filter["title"] != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(filter["title"].(string))+"%")
	}
	if filter["code"] != "" {
		query = query.Where("LOWER(code) LIKE ?", "%"+strings.ToLower(filter["code"].(string))+"%")
	}
	if filter["status"].(*int) != nil {
		query = query.Where("status = ?", *filter["status"].(*int))
	}

	if len(sort) > 0 {
		query = query.Order(sort)
	}

	var count int64
	pagination.Limit = limit
	pagination.Page = page
	err := query.Count(&count).Error
	if err != nil {
		return nil, nil, err
	}
	pagination.SetTotalRecords(count)
	err = query.Limit(pagination.Limit).Offset(pagination.GetOffset()).Find(&items).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return items, &pagination, nil
}

func (r *ShipmentPredefinedRepositoryImpl) GetShipmentPredefinedByUid(uid string) (*entity.ShippmentPredefined, error) {
	var db = r.base.GetDB()
	var ret *entity.ShippmentPredefined
	err := db.Where(&entity.ShippmentPredefined{BaseIDModel: base.BaseIDModel{UID: uid}}).Find(&ret).Error
	return ret, err
}

func (r *ShipmentPredefinedRepositoryImpl) UpdateShipmentPredefined(dto entity.ShippmentPredefined) (*entity.ShippmentPredefined, error) {
	var db = r.base.GetDB()
	var ret *entity.ShippmentPredefined

	ret, err := r.GetShipmentPredefinedByUid(dto.UID)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	if len(dto.Code) > 0 {
		data["code"] = dto.Code
	}
	if len(dto.Type) > 0 {
		data["type"] = dto.Type
	}
	if len(dto.Title) > 0 {
		data["title"] = dto.Title
	}
	if dto.Status != ret.Status {
		data["status"] = dto.Status
	}
	if dto.Note != ret.Note {
		data["note"] = dto.Note
	}

	err = db.Model(&entity.ShippmentPredefined{}).Where("uid=?", dto.UID).Updates(data).Error
	if err != nil {
		return nil, err
	}
	return r.GetShipmentPredefinedByUid(dto.UID)
}
