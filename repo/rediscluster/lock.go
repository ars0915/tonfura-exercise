package rediscluster

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func (r *RedisRepo) Lock(ctx context.Context, lockKey string) error {
	_, err := r.client.SetNX(ctx, lockKey, true, 3*time.Second).Result()
	if err != nil {
		return errors.Wrap(err, "set lock")
	}
	return nil
}

func (r *RedisRepo) UnLock(ctx context.Context, lockKey string) error {
	_, err := r.client.Del(ctx, lockKey).Result()
	if err != nil {
		return errors.Wrap(err, "unlock lock")
	}
	return nil
}
