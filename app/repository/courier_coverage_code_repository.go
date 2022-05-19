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
	CreateCourierCoverageCodeRepo(courierCoverageCode *entity.CourierCoverageCode) (*entity.CourierCoverageCode, error)
	GetCourierUid(courier *entity.Courier, uid string) error
	FindByParams(limit int, page int, sort string) ([]entity.CourierCoverageCode, *base.Pagination, error)
	FindByUid(uid string) (*entity.CourierCoverageCode, error)
	Update(uid string, input map[string]interface{}) error
}

func NewCourierCoverageCodeRepository(br BaseRepository) CourierCoverageCodeRepository {
	return &CourierCoverageCodeRepo{br}
}

func (r *CourierCoverageCodeRepo) CreateCourierCoverageCodeRepo(courierCoverageCode *entity.CourierCoverageCode) (*entity.CourierCoverageCode, error) {
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

func (r *CourierCoverageCodeRepo) FindByParams(limit int, page int, sort string) ([]entity.CourierCoverageCode, *base.Pagination, error) {
	var courierCoverageCodes []entity.CourierCoverageCode
	var pagination base.Pagination

	query := r.base.GetDB()

	if len(sort) > 0 {
		query = query.Order(sort)
	}

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.Paginate(courierCoverageCodes, &pagination, query, int64(len(courierCoverageCodes)))).
		Find(&courierCoverageCodes).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return courierCoverageCodes, &pagination, nil
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
		Where("uid=?", uid).
		First(&courierCoverageCode).Error
	if err != nil {
		return nil, err
	}

	return &courierCoverageCode, nil
}

func (r *CourierCoverageCodeRepo) Update(uid string, input map[string]interface{}) error {
	err := r.base.GetDB().Model(&entity.CourierCoverageCode{}).
		Where("uid=?", uid).
		Updates(input).Error
	if err != nil {
		return err
	}
	return nil
}
