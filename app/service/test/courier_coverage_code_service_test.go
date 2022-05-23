package test

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"os"
	"testing"
)

//var loggerCourierCoverageCode log.Logger

var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
	//db.AutoMigrate(&entity.CourierCoverageCode{})
}

func TestCreateCourierCoverageCode(t *testing.T) {
	req := request.SaveCourierCoverageCodeRequest{
		CourierUID:  "UCMvWngocMqKbaC3AWQBF",
		CountryCode: "VN",
		PostalCode:  "70000",
		Description: "Vietnam code",
	}

	result, _ := svcCourierCoverageCode.CreateCourierCoverageCode(req)
	assert.NotNil(t, result)
	assert.Equal(t, "UCMvWngocMqKbaC3AWQBF", result.CourierUID, "Courier UID is UCMvWngocMqKbaC3AWQBF")
	assert.Equal(t, "VN", result.CountryCode, "Courier UID is VN")
	assert.Equal(t, "70000", result.PostalCode, "Courier UID is 70000")
	assert.Equal(t, "Vietnam code", result.Description, "Courier UID is Vietnam code")

}
