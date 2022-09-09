package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v9"
	"k8s.io/klog/v2"
)

var ctx = context.Background()

// Prefix for all keys
const keyPrefix = "pippin"

// Singleton to keep assets loaded in memory
type redisManager struct {
	Client *redis.Client
	Locker *redislock.Client
	Mock   bool
}

var ErrLockNotObtained = errors.New("couldn't obtain lock")

// Retry every 100ms, for up-to 3x
var LockRetryStrategy = redislock.Options{
	RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 3),
}

var singleton *redisManager
var once sync.Once

func GetRedisDB() *redisManager {
	once.Do(func() {
		if utils.GetEnv("MOCK_REDIS", "false") == "true" {
			klog.Infof("Using mock redis client because MOCK_REDIS=true is set in environment")
			mr, _ := miniredis.Run()
			client := redis.NewClient(&redis.Options{
				Addr: mr.Addr(),
			})
			locker := redislock.New(client)
			singleton = &redisManager{
				Client: client,
				Locker: locker,
				Mock:   true,
			}
		} else {
			redis_port, err := strconv.Atoi(utils.GetEnv("REDIS_PORT", "6379"))
			if err != nil {
				panic("Invalid REDIS_PORT specified")
			}
			redis_db, err := strconv.Atoi(utils.GetEnv("REDIS_DB", "0"))
			if err != nil {
				panic("Invalid REDIS_DB specified")
			}
			client := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%d", utils.GetEnv("REDIS_HOST", "localhost"), redis_port),
				DB:   redis_db,
			})
			locker := redislock.New(client)
			singleton = &redisManager{
				Client: client,
				Locker: locker,
				Mock:   false,
			}
		}
	})
	return singleton
}

// del - Redis DEL
func (r *redisManager) Del(key string) (int64, error) {
	val, err := r.Client.Del(ctx, key).Result()
	return val, err
}

// get - Redis GET
func (r *redisManager) Get(key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	return val, err
}

// set - Redis SET
func (r *redisManager) Set(key string, value string, expiry time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiry).Err()
	return err
}

// hlen - Redis HLEN
func (r *redisManager) Hlen(key string) (int64, error) {
	val, err := r.Client.HLen(ctx, key).Result()
	return val, err
}

// hget - Redis HGET
func (r *redisManager) Hget(key string, field string) (string, error) {
	val, err := r.Client.HGet(ctx, key, field).Result()
	return val, err
}

// hgetall - Redis HGETALL
func (r *redisManager) Hgetall(key string) (map[string]string, error) {
	val, err := r.Client.HGetAll(ctx, key).Result()
	return val, err
}

// hset - Redis HSET
func (r *redisManager) Hset(key string, field string, values interface{}) error {
	err := r.Client.HSet(ctx, key, field, values).Err()
	return err
}

// hdel - Redis HDEL
func (r *redisManager) Hdel(key string, field string) error {
	err := r.Client.HDel(ctx, key, field).Err()
	return err
}
