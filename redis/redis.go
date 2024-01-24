package redis

import (
	"time"

	"github.com/picoorg/common/context"
)

type Redis interface {
	Get(ctx context.Context, key string) (result string, exist bool, err error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) (result []string, nextCursor uint64, exist bool, err error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error)
	GetDel(ctx context.Context, key string) (result string, exist bool, err error)
	Del(ctx context.Context, key string) (err error)
}
