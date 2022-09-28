package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/helper/global"
	"go-klikdokter/pkg/util"
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
	CombinationUnique(courierCoverageCode *entity.CourierCoverageCode, courierUid uint64, countryCode, postalCode, subdistrict string, id uint64) (int64, error)
	FindShipperCourierCoverage(input *request.FindShipperCourierCoverage) (*entity.CourierCoverageCode, error)
	FindInternalAndMerchantCourierCoverage(courierID []uint64, countryCode, postalCode string) map[string]bool
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

func (r *CourierCoverageCodeRepo) CombinationUnique(courierCoverageCode *entity.CourierCoverageCode, courierId uint64, countryCode, postalCode, subdistrict string, id uint64) (int64, error) {
	var result = r.base.GetDB()
	if id == 0 {
		result = result.First(&courierCoverageCode, "courier_id = ? AND country_code = ? AND postal_code = ? AND subdistrict = ?", courierId, countryCode, postalCode, subdistrict)
	} else {
		result = result.First(&courierCoverageCode, "courier_id = ? AND country_code = ? AND postal_code = ? AND subdistrict = ? AND id != ?", courierId, countryCode, postalCode, subdistrict, id)
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

	for k, v := range filters {

		if util.IsSliceAndNotEmpty(v) {

			switch k {
			case "courier_name":
				query = query.Joins("JOIN courier ON courier.id = courier_coverage_code.courier_id").
					Where(global.AddLike("courier.courier_name", v.([]string)))

			case "country_code", "postal_code", "subdistrict":
				query = query.Where(k+" IN ?", v.([]string))

			case "description", "":
				query = query.Where(global.AddLike("courier_coverage_code.description", v.([]string)))

			case "status":
				query = query.Where("courier_coverage_code.status IN ?", v)

			default:
				if string(k[0:4]) == "code" { //filtering for field code1, code2, .... code6
					query = query.Where(global.AddLike(k, v.([]string)))
				}
			}
		}
	}

	if len(sort) == 0 {
		sort = "courier_coverage_code.updated_at DESC"
	}

	query = query.Order(sort)
	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(&items, &pagination, query, int64(len(items)))).
		Find(&items).
		Error

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

func (r *CourierCoverageCodeRepo) FindShipperCourierCoverage(input *request.FindShipperCourierCoverage) (*entity.CourierCoverageCode, error) {
	var courierCoverageCode entity.CourierCoverageCode
	err := r.base.GetDB().
		Preload("Courier").
		Where("courier_id = ?", input.CourierID).
		Where("country_code = ?", input.CountryCode).
		Where("postal_code = ?", input.PostalCode).
		Where("subdistrict = ?", input.Subdistrict).
		First(&courierCoverageCode).Error
	if err != nil {
		return nil, err
	}
	if courierCoverageCode.Courier != nil {
		courierCoverageCode.CourierName = courierCoverageCode.Courier.CourierName
	}
	return &courierCoverageCode, nil
}

func (r *CourierCoverageCodeRepo) FindInternalAndMerchantCourierCoverage(courierID []uint64, countryCode, postalCode string) map[string]bool {
	var result = make(map[string]bool)
	var courierCoverageCode []entity.CourierCoverageCode
	var courierType = []string{"internal", "merchant"}

	err := r.base.GetDB().
		Preload("Courier").
		Joins("INNER JOIN courier c ON c.id = courier_coverage_code.courier_id").
		Where(&entity.CourierCoverageCode{CountryCode: countryCode}).
		Where(&entity.CourierCoverageCode{PostalCode: postalCode}).
		Where("courier_id IN ?", courierID).
		Where("c.courier_type IN ?", courierType).
		Find(&courierCoverageCode).Error

	if err != nil {
		return result
	}

	for _, v := range courierCoverageCode {
		result[v.Courier.Code] = true
	}

	return result
}
