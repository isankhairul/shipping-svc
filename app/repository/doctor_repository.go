package repository

import (
	"go-klikdokter/app/model/entity"
)

type doctorRepo struct {
	base BaseRepository
}

type DoctorRepository interface {
	FindByUid(uid *string) (*entity.Doctor, error)
	Create(product *entity.Doctor) (*entity.Doctor, error)
}

func NewDoctorRepository(br BaseRepository) DoctorRepository {
	return &doctorRepo{br}
}

func (r *doctorRepo) FindByUid(uid *string) (*entity.Doctor, error) {
	var doctor entity.Doctor
	err := r.base.GetDB().
		Where("uid=?", uid).
		First(&doctor).Error
	if err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *doctorRepo) Create(doctor *entity.Doctor) (*entity.Doctor, error) {
	err := r.base.GetDB().
		Create(doctor).Error
	if err != nil {
		return nil, err
	}

	return doctor, nil
}
