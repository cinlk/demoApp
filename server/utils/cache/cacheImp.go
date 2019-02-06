package cache

import (
	"demoApp/server/utils/errorStatus"
	"fmt"
	"github.com/pkg/errors"
	"goframework/cache"
	"net/http"
)

func GetKeyFromCache(key string) ([]byte, error) {

	res, err := cache.CacheProxy.Get(key)

	if err != nil {
		return nil, &errorStatus.AppError{
			Code: http.StatusNotFound,
			Err:  errors.Wrap(err, fmt.Sprintf("can't find key from cache ")),
		}
	}
	return res, err

}

func SetKeyInCache(key string, value []byte, timeout int) error {

	return cache.CacheProxy.SetEX(key, value, timeout)
}

func ClearTokenBy(key string) error {

	return cache.CacheProxy.Delete(key)
}
