package rds

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"runtime"
	"time"
)

func InitRedis(user, pwd, host, port string, db, cpu, min, max int) (*redis.Pool, error) {
	RedisPool := &redis.Pool{
		MaxIdle:     runtime.GOMAXPROCS(cpu) * 15,
		IdleTimeout: time.Duration(min) * time.Second,
		MaxActive:   max,
		Dial: func() (redis.Conn, error) {
			if len(user) > 0 && len(pwd) > 0 {
				c, err := redis.Dial(
					"tcp",
					host+":"+port,
					redis.DialDatabase(db),
					redis.DialUsername(user),
					redis.DialPassword(pwd),
					redis.DialKeepAlive(3*time.Second),
					redis.DialConnectTimeout(5*time.Second),
					redis.DialReadTimeout(2*time.Second),
					redis.DialWriteTimeout(3*time.Second),
				)
				if err != nil {
					return nil, err
				}
				return c, err
			} else {
				c, err := redis.Dial(
					"tcp",
					host+":"+port,
					redis.DialDatabase(db),
					redis.DialKeepAlive(3*time.Second),
					redis.DialConnectTimeout(5*time.Second),
					redis.DialReadTimeout(2*time.Second),
					redis.DialWriteTimeout(3*time.Second),
				)
				if err != nil {
					return nil, err
				}
				return c, err
			}
		},
	}

	return RedisPool, nil
}

type RedisLock struct {
	redisPool *redis.Pool
	resource  string
	expire    time.Duration
}

// 构造一个锁结构体
func NewRedisLock(pool *redis.Pool, resource string, expire time.Duration) *RedisLock {
	return &RedisLock{
		redisPool: pool,
		resource:  resource,
		expire:    expire,
	}
}

// 尝试获取锁
func (lock *RedisLock) TryLock() bool {
	conn := lock.redisPool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("SET", lock.resource, "1", "EX", int(lock.expire.Seconds()), "NX"))
	if err != nil {
		fmt.Println("尝试获取锁发生错误：", err)
		return false
	}

	return result == "OK"
}

// 释放锁
func (lock *RedisLock) Unlock() {
	conn := lock.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", lock.resource)
	if err != nil {
		fmt.Println("释放锁发生错误：", err)
	}
}

func SetRdsHMSetValue(r *redis.Pool, rdsKey string, ex any, v ...any) (err error) {
	conn := r.Get()
	defer conn.Close()
	//存入redis-hash
	args := redis.Args{}.Add(rdsKey).Add(v...)
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		RedisErrLog(rdsKey, "HMSET reids failed! err:", err)
		return
	}
	//设置过期时间
	_, err = conn.Do("EXPIRE", rdsKey, ex)
	if err != nil {
		RedisErrLog(rdsKey, " expire reids failed! err:[%v]", err)
		return
	}
	return
}
