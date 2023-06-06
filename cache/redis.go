package cache

import (
	"context"
	"log"
	"time"

	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client = nil
var ctx = context.Background()

func GetRedis() (*redis.Client, context.Context) {
	config, _ := initializers.LoadConfig("../")
	if rdb != nil {
		return rdb, ctx
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_URL,
		Password: config.REDIS_PASS, // password set
		DB:       0,                 // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis Init Failed")
	}
	return rdb, ctx
}

func SetValue(key, value string, time time.Duration) error {
	rdb, ctx := GetRedis()
	err := rdb.Set(ctx, key, value, time).Err()
	return err
}

func GetValue(key string) (string, error) {
	rdb, ctx := GetRedis()
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func DeleteValue(key string) error {
	rdb, ctx := GetRedis()
	err := rdb.Del(ctx, key).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}
