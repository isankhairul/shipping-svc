package service

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type DoctorService interface {
	CreateDoctor(input request.SaveDoctorRequest) (*entity.Doctor, message.Message)
	GetDoctor(uid string) (*entity.Doctor, message.Message)
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

// swagger:route POST /doctors/ Doctor SaveDoctorRequest
// Create Doctor
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *doctorServiceImpl) CreateDoctor(input request.SaveDoctorRequest) (*entity.Doctor, message.Message) {
	logger := log.With(s.logger, "ProductService", "CreateProduct")
	s.baseRepo.BeginTx()
	//Set request to entity
	doctor := entity.Doctor{
		Name:   input.Name,
		Gender: input.Gender,
	}

	result, err := s.doctorRepo.Create(&doctor)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseRepo.RollbackTx()
		return nil, message.ErrDB
	}
	s.baseRepo.CommitTx()

	return result, message.SuccessMsg
}

// swagger:route GET /doctors/{id} Get-Doctor doctor
// Get Doctor
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *doctorServiceImpl) GetDoctor(uid string) (*entity.Doctor, message.Message) {
	logger := log.With(s.logger, "ProductService", "GetProduct")

	result, err := s.doctorRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	if result == nil {
		return nil, message.ErrNoData
	}

	return result, message.SuccessMsg
}
