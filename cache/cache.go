// Package cache 缓存工具类，可以缓存各种类型包括 struct 对象
package cache

import (
	"encoding/json"
	"gohub/pkg/logger"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type CacheService struct {
	Store Store
}

var once sync.Once
var Cache *CacheService

func InitWithCacheStore(store Store) {
	once.Do(func() {
		Cache = &CacheService{
			Store: store,
		}
	})
}

func Exist(key string) (bool, error) {
	return Cache.Store.Exist(key)
}

func Set(key string, obj interface{}, expireTime time.Duration) {
	b, err := json.Marshal(&obj)
	logger.LogIf(err)
	Cache.Store.Set(key, string(b), expireTime)
}

func Get(key string) interface{} {
	return Cache.Store.Get(key)
}

func Has(key string) bool {
	return Cache.Store.Has(key)
}

func SetString(key string, val string, expireTime time.Duration) {
	Cache.Store.Set(key, val, expireTime)
}

// GetToObject 应该传地址，用法如下:
//
//	model := user.User{}
//	cache.GetObject("key", &model)
func GetToObject(key string, wanted interface{}) {
	val := Cache.Store.Get(key)
	if len(val) > 0 {
		err := json.Unmarshal([]byte(val), &wanted)
		logger.LogIf(err)
	}
}

func GetToString(key string) string {
	return cast.ToString(Get(key))
}

func GetToBool(key string) bool {
	return cast.ToBool(Get(key))
}

func GetToInt(key string) int {
	return cast.ToInt(Get(key))
}

func GetToInt32(key string) int32 {
	return cast.ToInt32(Get(key))
}

func GetToInt64(key string) int64 {
	return cast.ToInt64(Get(key))
}

func GetToUint(key string) uint {
	return cast.ToUint(Get(key))
}

func GetToUint32(key string) uint32 {
	return cast.ToUint32(Get(key))
}

func GetToUint64(key string) uint64 {
	return cast.ToUint64(Get(key))
}

func GetToFloat64(key string) float64 {
	return cast.ToFloat64(Get(key))
}

func GetToTime(key string) time.Time {
	return cast.ToTime(Get(key))
}

func GetToDuration(key string) time.Duration {
	return cast.ToDuration(Get(key))
}

func GetToIntSlice(key string) []int {
	return cast.ToIntSlice(Get(key))
}

func GetToStringSlice(key string) []string {
	return cast.ToStringSlice(Get(key))
}

func GetToStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(Get(key))
}

func GetToStringMapString(key string) map[string]string {
	return cast.ToStringMapString(Get(key))
}

func GetToStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(Get(key))
}

func Forget(key string) {
	Cache.Store.Forget(key)
}

func Forever(key string, value string) {
	Cache.Store.Set(key, value, 0)
}

func Flush() {
	Cache.Store.Flush()
}

func PreDelAll(pre string) error {
	return Cache.Store.PreDelAll(pre)
}

func Delete(key string) bool {
	return Cache.Store.Delete(key)
}

func Increment(parameters ...interface{}) {
	Cache.Store.Increment(parameters...)
}

func Decrement(parameters ...interface{}) {
	Cache.Store.Decrement(parameters...)
}

func IsAlive() error {
	return Cache.Store.IsAlive()
}
