package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

var Redis = &redisClient{}

type redisClient struct {
	Prefix string
	*redis.Client
}

func InitRedis() {
	config := Configs.Redis
	r := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		//TODO other config
	})
	Redis.Client = r
	Redis.Prefix = "chat:"
}

func (r *redisClient) Set(c *gin.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(c, r.Prefix+key, value, expiration).Err()
	if err != nil {
		Logger.WithContext(c).Error(err)
		return err
	}
	return nil
}

func (r *redisClient) Get(c *gin.Context, key string) *redis.StringCmd {
	return r.Client.Get(c, r.Prefix+key)
}

func (r *redisClient) Lock(c *gin.Context, key string, expiration time.Duration) error {
	res, err := r.Client.SetNX(c, r.Prefix+key, "1", expiration).Result()
	if err != nil {
		return err
	}
	if res == false {
		return fmt.Errorf("获取锁失败")
	}
	return nil //获取锁成功
}

func (r *redisClient) Del(c *gin.Context, keys ...string) (delCount int64, err error) {
	newKeys := make([]string, 0)
	for _, key := range keys {
		newKeys = append(newKeys, r.Prefix+key)
	}
	delCount, err = r.Client.Del(c, newKeys...).Result()
	if err != nil {
		return
	}
	return
}
