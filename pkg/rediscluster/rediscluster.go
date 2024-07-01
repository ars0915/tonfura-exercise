package rediscluster

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/ars0915/tonfura-exercise/config"
)

func NewRedisClient(config config.ConfENV) (*redis.ClusterClient, error) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: strings.Split(config.Redis.Hosts, ","),
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "redis connect failed")
	}

	return redisClient, nil
}
