package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"course_online_backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var ctx = context.Background()

func RedisConn() *redis.Client {
	if RDB != nil {
		return RDB
	}

	conf := config.DbConfig

	dbNum, err := strconv.Atoi(conf.REDISDb)
	if err != nil {
		log.Println("Invalid REDIS_DB, using 0 as default")
		dbNum = 0
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.REDISHost, conf.REDISPort),
		Password: conf.REDISPassword, 
		DB:       dbNum,
	})

	pong, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Redis connected:", pong)
	}

	return RDB
}

func Set(key string, value interface{}, ttl time.Duration) error {
	return RDB.Set(ctx, key, value, ttl).Err()
}

func Get(key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}
