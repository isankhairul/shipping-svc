package service

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
)

type DoctorService interface {
	CreateDoctor(input request.SaveDoctorRequest) (*entity.Doctor, int, string)
	GetDoctor(uid string) (*entity.Doctor, int, string)
}

type doctorServiceImpl struct {
	logger     log.Logger
	baseRepo   repository.BaseRepository
	doctorRepo repository.DoctorRepository
}

func NewDoctorService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.DoctorRepository,
) DoctorService {
	return &doctorServiceImpl{lg, br, pr}
}

func (s *doctorServiceImpl) CreateDoctor(input request.SaveDoctorRequest) (*entity.Doctor, int, string) {
	logger := log.With(s.logger, "ProductService", "CreateProduct")
	s.baseRepo.BeginTx()
	//Set request to entity
	doctor := entity.Doctor{
		Name:   input.Name,
		Gender: input.Gender,
	}

	result, err := s.doctorRepo.Create(&doctor)
	if err != nil {
		level.Error(logger).Log(err)
		s.baseRepo.RollbackTx()
		return nil, message.CODE_ERR_DB, message.MSG_ERR_SAVE_DATA
	}
	s.baseRepo.CommitTx()

	return result, message.CODE_SUCCESS, ""
}

func (s *doctorServiceImpl) GetDoctor(uid string) (*entity.Doctor, int, string) {
	logger := log.With(s.logger, "ProductService", "GetProduct")

	result, err := s.doctorRepo.FindByUid(&uid)
	if err != nil {
		level.Error(logger).Log(err)
		return nil, message.CODE_ERR_DB, message.MSG_ERR_DB
	}

	if result == nil {
		return nil, message.CODE_ERR_DB, message.MSG_NO_DATA
	}

	return result, message.CODE_SUCCESS, ""
}
