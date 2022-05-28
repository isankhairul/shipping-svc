package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ChannelRepositoryMock struct {
	Mock mock.Mock
}

func (repository *ChannelRepositoryMock) FindByUid(uid *string) (*entity.Channel, error) {
	arguments := repository.Mock.Called(uid)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		channel := arguments.Get(0).(entity.Channel)
		return &channel, nil
	}
}

func (repository *ChannelRepositoryMock) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Channel, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort, filter)
	return arguments.Get(0).([]entity.Channel), arguments.Get(1).(*base.Pagination), nil
}

func (repository *ChannelRepositoryMock) FindAll(limit int, page int, sort string) ([]entity.Channel, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort)
	return arguments.Get(0).([]entity.Channel), arguments.Get(1).(*base.Pagination), nil
}

func (repository *ChannelRepositoryMock) CreateChannel(channel *entity.Channel) (*entity.Channel, error) {
	return channel, nil
}

func (repository *ChannelRepositoryMock) Delete(uid string) error {
	return nil
}

func (repository *ChannelRepositoryMock) Update(uid string, input map[string]interface{}) error {
	return nil
}

func (repository *ChannelRepositoryMock) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	return nil
}

func (repository *ChannelRepositoryMock) CheckExistsByUIdChannelCode(uid, channelCode string) (bool, error) {
	return false, nil
}
