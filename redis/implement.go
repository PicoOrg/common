package redis

import (
	"errors"
	"time"

	"github.com/picoorg/common/context"
	"github.com/redis/go-redis/v9"
)

func New(addr, password string, db int) Redis {
	return &implement{
		redis: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

type implement struct {
	redis *redis.Client
}

func (m *implement) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := m.redis.Set(ctx.GetRawContext(), key, value, expiration).Err()
	if err != nil {
		ctx.LogWithField("error", err).Error("failed to set redis")
		return err
	}
	return nil
}

func (m *implement) Get(ctx context.Context, key string) (string, bool, error) {
	result, err := m.redis.Get(ctx.GetRawContext(), key).Result()
	if errors.Is(err, redis.Nil) {
		ctx.LogWithField("error", err).Info("failed to get redis")
		return "", false, nil
	} else if err != nil {
		ctx.LogWithField("error", err).Error("failed to get redis")
		return "", false, err
	}
	return result, true, nil
}

func (m *implement) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, bool, error) {
	result, nextCursor, err := m.redis.Scan(ctx.GetRawContext(), cursor, match, count).Result()
	if errors.Is(err, redis.Nil) {
		ctx.LogWithField("error", err).Info("failed to get redis")
		return nil, 0, false, nil
	} else if err != nil {
		ctx.LogWithField("error", err).Error("failed to get redis")
		return nil, 0, false, err
	}
	return result, nextCursor, true, nil
}

func (m *implement) Del(ctx context.Context, key string) error {
	err := m.redis.Del(ctx.GetRawContext(), key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		ctx.LogWithField("error", err).Error("failed to del redis")
		return err
	}
	return nil
}
func (m *implement) GetDel(ctx context.Context, key string) (string, bool, error) {
	results, _ := m.redis.Pipelined(ctx.GetRawContext(), func(pipe redis.Pipeliner) error {
		pipe.Get(ctx.GetRawContext(), key).Result()
		pipe.Del(ctx.GetRawContext(), key).Err()
		return nil
	})
	result, err := results[0].(*redis.StringCmd).Result()
	if errors.Is(err, redis.Nil) {
		ctx.LogWithField("error", err).Error("failed to get del redis")
		return "", false, nil
	} else if err != nil {
		ctx.LogWithField("error", err).Error("failed to get del redis")
		return "", false, err
	}
	err = results[1].Err()
	if err != nil {
		ctx.LogWithField("error", err).Error("failed to get redis")
		return "", false, err
	}
	return result, true, nil
}
