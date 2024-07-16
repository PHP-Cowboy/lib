// Package jwt 处理 JWT 认证
package jwt

import (
	"errors"
	"fmt"
	jwtpkg "github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type JWTConfig struct {
	MaxTime time.Duration
	Key     []byte
}

type JWTData struct {
	Uid          uint64 `json:"uid"`
	UserName     string `json:"UserName"`
	Phone        string `json:"Phone"`
	Email        string `json:"Email"`
	Icon         int8   `json:"Icon"`
	RegTimestamp int64  `json:"reg_timestamp"`
	ChannelId    int    `json:"channel_id"`
	jwtpkg.StandardClaims
}

// 解析token
func (config *JWTConfig) GetUserInfoByToken(token string) (*JWTData, error) {
	var userInfo JWTData
	data, err := jwtpkg.ParseWithClaims(token, &userInfo, func(token *jwtpkg.Token) (interface{}, error) {
		return config.Key, nil
	})

	if data != nil && data.Valid {
		return &userInfo, nil
	} else if ve, ok := err.(*jwtpkg.ValidationError); ok {
		if ve.Errors&jwtpkg.ValidationErrorMalformed != 0 {
			return nil, errors.New("不是一个合法的token")
		} else if ve.Errors&(jwtpkg.ValidationErrorExpired|jwtpkg.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("token过期了")
		} else if ve.Errors&jwtpkg.ValidationErrorIssuedAt != 0 {
			curTime := time.Now()
			if userInfo.IssuedAt > curTime.Add(-5*time.Minute).Unix() && userInfo.IssuedAt < curTime.Unix() {
				return &userInfo, nil
			} else {
				return nil, errors.New("Token iat error")
			}
		} else {
			return nil, errors.New("无法处理这个token")
		}
	} else {
		return nil, errors.New("无法处理这个token")
	}
}

func (config *JWTConfig) GetTokenByUserInfo(user JWTData) (string, error) {
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, user)
	signToken, err := token.SignedString(config.Key)

	if err != nil {
		return "", err
	}

	return signToken, nil
}

// redis key
const (
	RedisKeyTokenId = "jwt:"
)

// 暂时都是调用这个接口，后续加个缓存 todo
func GetTokenByUid(r *redis.Pool, id uint64) (string, error) {
	if r == nil {
		return "", errors.New(fmt.Sprintf("GetTokenByUid: redis err:", id))
	}

	rdsConn := r.Get()
	defer rdsConn.Close()
	rdsKey := RedisKeyTokenId + strconv.FormatUint(id, 10)
	token, err := redis.String(rdsConn.Do("Get", rdsKey))
	if err != nil {
		return "", errors.New(fmt.Sprintf("GetTokenByUid: get redis failed! err:", err, id))

	}
	return token, err
}

//func  GetTokenAndTTLByUid(id uint64) (token string, expire time.Duration, err error) {
//	rdsConn := config.RedisPool.Get()
//	defer rdsConn.Close()
//	rdsKey := jwt.RedisKeyTokenId + strconv.FormatUint(id, 10)
//	res, err := redis.String(rdsConn.Do("Get", rdsKey))
//	// 如果键存在并且有设置过期时间，获取TTL
//	ttlResult, err := redis.Int(rdsConn.Do("TTL", rdsKey))
//	if ttlResult != -1 {
//		// ttlResult是剩余秒数，可以进行相关处理
//	} else if err != nil {
//		// 处理TTL命令的错误
//	} else {
//		// 键不存在或者没有设置过期时间
//	}
//
//	return
//}
