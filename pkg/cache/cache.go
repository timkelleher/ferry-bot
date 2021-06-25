package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
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

func Get(ctx context.Context, key string, tty time.Duration, getter func(args wsdot.EndpointArguments) (string, error), args wsdot.EndpointArguments) (string, error) {
	redisCache := cache.New(&cache.Options{
		Redis:      getConn(),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	// todo: make key unique per request w/ args
	//key += sha256.Sum256(args)

	var wanted string
	if cacheOn {
		if err := redisCache.Get(ctx, key, &wanted); err == nil {
			fmt.Println("Fetch cached value for key", key)
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
			fmt.Printf("Cache Error with key %s: %s", key, err.Error())
			return data, nil
		}
	}
	return data, nil
}
