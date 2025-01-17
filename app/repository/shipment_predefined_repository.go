package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/pkg/util"

	"gorm.io/gorm"
)

type ShipmentPredefinedRepository interface {
	GetAll(limit int, page int, sort string, filter map[string]interface{}) ([]*entity.ShippmentPredefined, *base.Pagination, error)
	GetListByType(Type string) ([]entity.ShippmentPredefined, error)
	UpdateShipmentPredefined(dto entity.ShippmentPredefined) (*entity.ShippmentPredefined, error)
	GetShipmentPredefinedByUid(uid string) (*entity.ShippmentPredefined, error)
}

type ShipmentPredefinedRepositoryImpl struct {
	base BaseRepository
}

func NewShipmentPredefinedRepository(br BaseRepository) ShipmentPredefinedRepository {
	return &ShipmentPredefinedRepositoryImpl{br}
}

func (r *ShipmentPredefinedRepositoryImpl) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)
	pagination.SetTotalRecords(totalRecords)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

// GetAll implements ShipmentPredefinedRepository
func (r *ShipmentPredefinedRepositoryImpl) GetAll(limit int, page int, sort string, filter map[string]interface{}) ([]*entity.ShippmentPredefined, *base.Pagination, error) {
	var items []*entity.ShippmentPredefined
	var pagination base.Pagination

	query := r.base.GetDB().Model(&entity.ShippmentPredefined{})

	for k, v := range filter {

		if util.IsSliceAndNotEmpty(v) {

			switch k {
			case "note", "title", "code":
				query = query.Where(like(k, v.([]string)))

			case "type":
				query = query.Where(like(k, v.([]string)))

			case "status":
				query = query.Where("status IN ?", v)

			}
		}
	}

	if len(sort) == 0 {
		sort = "updated_at DESC"
	}

	query = query.Order(sort)

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(items, &pagination, query, int64(len(items)))).
		Find(&items).
		Error

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

func (r *ShipmentPredefinedRepositoryImpl) GetListByType(Type string) ([]entity.ShippmentPredefined, error) {
	var result []entity.ShippmentPredefined

	err := r.base.GetDB().Model(&entity.ShippmentPredefined{}).
		Where(&entity.ShippmentPredefined{Type: Type}).
		Order("id desc").
		Find(&result).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}
