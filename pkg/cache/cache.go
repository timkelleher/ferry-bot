package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/timkelleher/ferry-bot/pkg/wsdot"
)

//todo: make this an env file
var cacheOn = true
var conn *redis.Ring

func getConn() *redis.Ring {
	if conn == nil {
		conn = redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				"server1": ":6379",
			},
		})
	}
	return conn
}

func Get(key string, tty time.Duration, logger *logrus.Logger, getter func(args wsdot.EndpointArguments) (string, error), args wsdot.EndpointArguments) (string, error) {
	ctx := context.Background()

	redisCache := cache.New(&cache.Options{
		Redis:      getConn(),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	// todo: make key unique per request w/ args
	//key += sha256.Sum256(args)

	var wanted string
	if cacheOn {
		if err := redisCache.Get(ctx, key, &wanted); err == nil {
			logger.WithField("key", key).Info("Fetching cached value")
			return wanted, nil
		}
	}

	data, err := getter(args)
	if err != nil {
		return "", err
	}
	if cacheOn {
		if err := redisCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: data,
			TTL:   time.Hour,
		}); err != nil {
			logger.WithFields(logrus.Fields{
				"key":   key,
				"error": err.Error(),
			}).Error("Fetching cached value")
			return data, nil
		}
		logger.WithField("key", key).Info("Caching value")
	}
	return data, nil
}
