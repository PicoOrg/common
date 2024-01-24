package eip

import (
	"fmt"

	"github.com/picoorg/common/context"
)

func (m *implement) getEdges(ctx context.Context, cityhash string, num int) ([]Edge, error) {
	for page := 0; page < 100; page += 1 {
		edges := make([]Edge, 0)
		for cursor, left := uint64(0), int64(num); ; {
			keys, nextCursor, _, err := m.redis.Scan(ctx, cursor, m.getEdgeNewAllRedisKey(), left)
			if err != nil {
				ctx.LogWithField("err", err).Error("redis scan error")
				return nil, err
			}

			for _, key := range keys {
				edge, err := NewEdgeFromRedisKey(ctx, m.redis, key)
				if err != nil {
					return nil, err
				}
				edges = append(edges, edge)
			}

			cursor = nextCursor
			left = int64(num) - int64(len(edges))

			if left <= 0 {
				return edges, nil
			} else if cursor == 0 {
				break
			}
		}

		edges, err := m.EdgeDevice(ctx, cityhash, page*1000, 1000)
		if err != nil {
			return nil, err
		}
		for _, edge := range edges {
			isNew, err := edge.IsNew(ctx)
			if err != nil {
				return nil, err
			} else if isNew {
				err = edge.SetNew(ctx)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	err := fmt.Errorf("get edges error")
	ctx.LogWithField("err", err).LogWithField("cityhash", cityhash).LogWithField("num", num).Error("get edges error")
	return nil, err
}
