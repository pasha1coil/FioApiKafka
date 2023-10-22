package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func InitRedis() (*redis.Client, error) {
	cache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Infoln("Check redis...")

	log.Infoln("Set. Key - Ping, Value - Pong")
	err := cache.Set(ctx, "Ping", "Pong", 0).Err()
	if err != nil {
		log.Errorf("Error test set to redis:%s", err.Error())
		return nil, err
	}

	log.Infoln("Get. Key - Ping")
	value, err := cache.Get(ctx, "Ping").Result()
	if err != nil {
		log.Errorf("Error test get from redis:%s", err.Error())
		return nil, err
	}
	log.Infof("Retrieved value from redis:%s", value)

	log.Infoln("Delete. Key - Ping")
	_, err = cache.Del(ctx, "Ping").Result()
	if err != nil {
		log.Errorf("Error test deleted from redis:%s", err.Error())
		return nil, err
	}
	return cache, nil
}
