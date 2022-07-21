package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"math"

	"gorm.io/gorm"
)

type CourierCoverageCodeRepo struct {
	base BaseRepository
}

type CourierCoverageCodeRepository interface {
	Create(courierCoverageCode *entity.CourierCoverageCode) (*entity.CourierCoverageCode, error)
	GetCourierUid(courier *entity.Courier, uid string) error
	GetCourierId(courier *entity.Courier, id uint64) error
	FindByParams(limit int, page int, sort string, filters map[string]interface{}) ([]*entity.CourierCoverageCode, *base.Pagination, error)
	FindByUid(uid string) (*entity.CourierCoverageCode, error)
	Update(uid string, values map[string]interface{}) (*entity.CourierCoverageCode, error)
	CombinationUnique(courierCoverageCode *entity.CourierCoverageCode, courierUid uint64, countryCode, postalCode string, id uint64) (int64, error)
	DeleteByUid(uid string) error
}

func NewCourierCoverageCodeRepository(br BaseRepository) CourierCoverageCodeRepository {
	return &CourierCoverageCodeRepo{br}
}

func (r *CourierCoverageCodeRepo) Create(courierCoverageCode *entity.CourierCoverageCode) (*entity.CourierCoverageCode, error) {
	err := r.base.GetDB().Create(courierCoverageCode).Error

	if err != nil {
		return nil, err
	}

	return courierCoverageCode, nil

}

func (r *CourierCoverageCodeRepo) GetCourierUid(courier *entity.Courier, uid string) error {
	err := r.base.GetDB().First(courier, "uid = ?", uid).Error

	if err != nil {
		return err
	}

	return nil
}
func (r *CourierCoverageCodeRepo) GetCourierId(courier *entity.Courier, id uint64) error {
	err := r.base.GetDB().First(courier, "id = ?", id).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *CourierCoverageCodeRepo) CombinationUnique(courierCoverageCode *entity.CourierCoverageCode, courierId uint64, countryCode, postalCode string, id uint64) (int64, error) {
	var result = r.base.GetDB()
	if id == 0 {
		result = result.First(&courierCoverageCode, "courier_id = ? AND country_code = ? AND postal_code = ?", courierId, countryCode, postalCode)
	} else {
		result = result.First(&courierCoverageCode, "courier_id = ? AND country_code = ? AND postal_code = ? AND id != ?", courierId, countryCode, postalCode, id)
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, result.Error
	}

	return result.RowsAffected, nil

}

func (r *CourierCoverageCodeRepo) FindByParams(limit int, page int, sort string, filters map[string]interface{}) ([]*entity.CourierCoverageCode, *base.Pagination, error) {
	var items []*entity.CourierCoverageCode
	var pagination base.Pagination

	db := r.base.GetDB()
	query := db.Model(entity.CourierCoverageCode{}).Preload("Courier").Joins("Courier")

	if filters["courier_name"] != "" {
		query = query.Joins("JOIN courier ON courier.id = courier_coverage_code.courier_id AND courier.courier_name = ?", filters["courier_name"].(string))
	}
	if filters["country_code"] != "" {
		query = query.Where(entity.CourierCoverageCode{CountryCode: filters["country_code"].(string)})
	}
	if filters["postal_code"] != "" {
		query = query.Where(entity.CourierCoverageCode{PostalCode: filters["postal_code"].(string)})
	}
	if filters["description"] != "" {
		query = query.Where(entity.CourierCoverageCode{Description: filters["description"].(string)})
	}

	if filters["status"].(*int) != nil {
		query = query.Where(entity.CourierCoverageCode{Status: *filters["status"].(*int)})
	}
	if len(sort) > 0 {
		query = query.Order(sort)
	} else {
		query = query.Order("courier_coverage_code.updated_at DESC")
	}

	pagination.Limit = limit
	pagination.Page = page

	var totalRecords int64
	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, nil, err
	}

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	// err = query.Scopes(r.Paginate(&pagination)).Find(&items).Error
	err = query.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Find(&items).Error
	if err != nil {
		return nil, nil, err
	}
	for _, item := range items {
		if item.Courier != nil {
			item.CourierName = item.Courier.CourierName
			item.CourierUID = item.Courier.UID
		}
	}
	return items, &pagination, nil
}

func (r *CourierCoverageCodeRepo) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(currRecord)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *CourierCoverageCodeRepo) FindByUid(uid string) (*entity.CourierCoverageCode, error) {
	var courierCoverageCode entity.CourierCoverageCode
	err := r.base.GetDB().
		Preload("Courier").
		Where("uid = ?", uid).
		First(&courierCoverageCode).Error
	if err != nil {
		return nil, err
	}
	if courierCoverageCode.Courier != nil {
		courierCoverageCode.CourierName = courierCoverageCode.Courier.CourierName
	}
	return &courierCoverageCode, nil
}

func (r *CourierCoverageCodeRepo) Update(uid string, values map[string]interface{}) (*entity.CourierCoverageCode, error) {
	var courierCoverageCode entity.CourierCoverageCode
	err := r.base.GetDB().Model(&courierCoverageCode).
		Where("uid=?", uid).
		Updates(values).Error
	if err != nil {
		return nil, err
	}
	return &courierCoverageCode, nil
}

func (r *CourierCoverageCodeRepo) DeleteByUid(uid string) error {

	var ret entity.CourierCoverageCode
	err := r.base.GetDB().Where("uid=?", uid).First(&ret).Error

	if err != nil {
		return err
	}

	err = r.base.GetDB().Where("uid=?", uid).Delete(&ret).Error

	if err != nil {
		return err
	}

	return nil
}
