package im_services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

func throttleKey(key string) string {
	return fmt.Sprintf("throttle:%s", key)
}

func ThrottleWithinLimit(rdb *redis.Client, ctx context.Context, key string, limit time.Duration, expires_in time.Duration) (bool, error) {
	limitSecond := int64(limit.Seconds())
	fmt.Println("limitSecond")
	fmt.Println(limitSecond)
	b, err := rdb.Incr(ctx, throttleKey(key)).Result()
	if err != nil {
		return false, err
	}
	fmt.Println("b")
	fmt.Println(b)

	if b == 1 {
		if _, err := rdb.Expire(ctx, throttleKey(key), expires_in).Result(); err != nil {
			return false, err
		}
	}

	if b <= limitSecond {
		return true, nil
	}
	return false, nil
}
