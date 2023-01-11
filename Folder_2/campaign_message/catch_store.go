package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type redisClient struct {
	*redis.Client
}

func openRedisStore(host string, port int, password string, db int) (*redisClient, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // if pass not set, use empty is ""
		DB:       db,       // default DB is 0
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, &RedisError{
			Type: "REDIS",
			Err:  fmt.Errorf("PING: %v", err),
		}
	}
	return &redisClient{
		rdb,
	}, nil
}

func (rdb *redisClient) setKey(key string, value string) {
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (rdb *redisClient) getDailyCap(messageTemplateId string, playerId string) string {
	capKeys := getCapKeys(messageTemplateId, playerId)

	val, err := rdb.Get(ctx, capKeys["daily"]).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func (rdb *redisClient) getMonthlyCap(messageTemplateId string, playerId string) string {
	capKeys := getCapKeys(messageTemplateId, playerId)

	val, err := rdb.Get(ctx, capKeys["monthly"]).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func (rdb *redisClient) getTotalCap(messageTemplateId string, playerId string) string {
	capKeys := getCapKeys(messageTemplateId, playerId)

	val, err := rdb.Get(ctx, capKeys["total"]).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func getCapKeys(messageTemplateId string, playerId string) map[string]string {
	key := fmt.Sprintf("mt/%s/pl/%s", messageTemplateId, playerId)
	// key = "mt/#{id}/pl/#{player_id}"
	capKeys := map[string]string{
		"total":   key,
		"daily":   key + "/daily",
		"monthly": key + "/monthly",
	}

	return capKeys
}

func (rdb *redisClient) IncrementCaps(messageTemplateId string, playerId string) error {
	capKeys := getCapKeys(messageTemplateId, playerId)

	rdb.Incr(ctx, capKeys["total"])
	rdb.Incr(ctx, capKeys["daily"])
	rdb.Incr(ctx, capKeys["monthly"])

	if _, err := rdb.Expire(ctx, capKeys["daily"], time.Hour*time.Duration(12)).Result(); err != nil {
		return err
	}

	if _, err := rdb.Expire(ctx, capKeys["monthly"], time.Hour*time.Duration(24*30)).Result(); err != nil {
		return err
	}

	return nil
}

// this is sample data. it will be removed in production
func (rdb *redisClient) setSampleCapsValue(messageTemplateId string, playerId string) {

	capKeys := getCapKeys(messageTemplateId, playerId)

	for _, keyVal := range capKeys {
		err := rdb.Set(ctx, keyVal, 1, 0).Err()
		if err != nil {
			panic(err)
		}
	}

}
