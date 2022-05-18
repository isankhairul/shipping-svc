package repository

import "go-klikdokter/app/model/entity"

type CourierCoverageCodeRepo struct {
	base BaseRepository
}

type CourierCoverageCodeRepository interface {
	CreateCoverageCodeRepo(courierCoverageCode *entity.CourierCoverageCode) (*entity.CourierCoverageCode, error)
}

func NewCourierCoverageCodeRepository(br BaseRepository) CourierCoverageCodeRepository {
	return &CourierCoverageCodeRepo{br}
}

func (r *CourierCoverageCodeRepo) CreateCoverageCodeRepo(courierCoverageCode *entity.CourierCoverageCode) (*entity.CourierCoverageCode, error) {
	err := r.base.GetDB().Create(courierCoverageCode).Error

	if err != nil {
		return nil, err
	}

	return courierCoverageCode, nil

}
