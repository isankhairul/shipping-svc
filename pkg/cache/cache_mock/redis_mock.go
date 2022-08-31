package cache_mock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type Redis_Mock struct {
	Mock mock.Mock
}

func (cache *Redis_Mock) Deletes(keys ...string) {
}

func (cache *Redis_Mock) Delete(key string) {
}

func (cache *Redis_Mock) HashDelete(key string, field string) {
}

func (cache *Redis_Mock) HashDeletes(key string, fields ...string) {
}

func (cache *Redis_Mock) DeletePrefix(prefix string) {
}

func (cache *Redis_Mock) Expire(key string, expiration time.Duration) {
}

func (cache *Redis_Mock) Get(key string) (string, error) {
	arguments := cache.Mock.Called()

	if len(arguments) > 1 && arguments.Get(1) != nil {
		return "", arguments.Get(1).(error)
	}

	if arguments.Get(0) == nil {
		return "", nil
	}

	return arguments.Get(0).(string), nil
}

func (cache *Redis_Mock) GetJsonStruct(key string, structObj interface{}) error {
	arguments := cache.Mock.Called()

	if arguments.Get(0) == nil {
		return nil
	}

	return arguments.Get(0).(error)
}

func (cache *Redis_Mock) HashGet(key string, field string) (string, error) {
	arguments := cache.Mock.Called()

	if len(arguments) > 1 && arguments.Get(1) != nil {
		return "", arguments.Get(1).(error)
	}

	if arguments.Get(0) == nil {
		return "", nil
	}

	return arguments.Get(0).(string), nil
}

func (cache *Redis_Mock) HashGetJsonStruct(key string, field string, structObj interface{}) error {
	arguments := cache.Mock.Called()

	if arguments.Get(0) == nil {
		return nil
	}

	return arguments.Get(0).(error)
}

func (cache *Redis_Mock) IsExist(key string) bool {
	arguments := cache.Mock.Called()

	if arguments.Get(0) == nil {
		return false
	}

	return arguments.Get(0).(bool)
}

func (cache *Redis_Mock) Set(key string, value interface{}, exp ...int) {
}

func (cache *Redis_Mock) HashSet(key string, field string, value interface{}, exp ...int) error {
	arguments := cache.Mock.Called()

	if arguments.Get(0) == nil {
		return nil
	}

	return arguments.Get(0).(error)
}

func (cache *Redis_Mock) SetJsonStruct(key string, value interface{}, exp ...int) {
}
