package tool

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"strconv"
)

var (
	TimeFormatYmdHis = "2006-01-02 15:04:05"
	TimeFormatYmd    = "20060102"
	TimeFormatYm     = "200601"
	TimeFormatYmdH   = "2006010215"
)

// md5 加密
func GetMd5(str string) (code string) {
	md := md5.New()
	md.Write([]byte(str))
	code = hex.EncodeToString(md.Sum(nil))
	return
}

// 格式化json字符串
func JsonString(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 字符串拼接
func StringJoin(str []string) string {
	var st bytes.Buffer
	for _, val := range str {
		st.WriteString(val)
	}
	return st.String()
}
func IntJoin2String(intS []int, sep string) string {
	var st bytes.Buffer
	var isFirst = true
	for _, val := range intS {
		if len(sep) > 0 && !isFirst {
			st.WriteString(sep)
		}
		st.WriteString(strconv.Itoa(val))

		isFirst = false
	}
	return st.String()
}

// 计算分页
func GetPage(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 30
	}

	limit = pageSize
	offset = (page - 1) * pageSize
	return
}

// slice 转 map
func SliceToMap[T comparable](arr []T) (mp map[T]struct{}) {
	mp = make(map[T]struct{}, 0)
	for _, t := range arr {
		mp[t] = struct{}{}
	}
	return mp
}

// map 转 slice 不去重
func MapToSlice[T comparable](mp map[T]struct{}) (arr []T) {
	for id, _ := range mp {
		arr = append(arr, id)
	}
	return arr
}
func MarkFlag(src, flag int) (dest int) {
	dest = src | flag
	return
}
func UnmarkFlag(src, flag int) (dest int) {
	dest = src &^ flag
	return
}
func IsMarkFlag(src, flag int) bool {
	if (src & flag) != 0 {
		return true
	}
	return false
}

// 传入一个权重数组，获得下标
func GetRandFromWeight(weights []int32) int32 {
	// 计算权重的总和
	totalWeight := int32(0)
	for _, w := range weights {
		totalWeight += w
	}
	if totalWeight == 0 { //总权重不能为0
		return -1
	}

	// 生成随机数
	r := rand.Int31n(totalWeight)

	// 根据随机数获取对应值
	sel := int32(0)
	sum := int32(0)
	for i, w := range weights {
		sum += w
		if r < sum {
			sel = int32(i)
			break
		}
	}
	return sel
}

// 获得redis过期时间，2H+1H随机数
func GetRedisExpireTime() int {
	expire := 2 * 3600
	return expire + rand.Intn(3600)
}

// 把数字转为58进制
func Numberto58Base(number uint64) string {
	digits := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

	var result string
	for number > 0 {
		remainder := number % 58
		result = string(digits[remainder]) + result
		number /= 58
	}
	return result
}
