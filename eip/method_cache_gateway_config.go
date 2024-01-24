package eip

import (
	"encoding/json"

	"github.com/picoorg/common/context"
)

func (m *implement) CacheGatewayConfig(ctx context.Context, c *gatewayConfig) error {
	err := m.redis.Set(ctx, m.getConfigIDRedisKey(), c.ID, 0)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&c.Rules)
	if err != nil {
		ctx.LogWithField("err", err).Error("json marshal error")
		return err
	}

	err = m.redis.Set(ctx, m.getConfigRulesRedisKey(), string(data), 0)
	if err != nil {
		return err
	}

	for _, t := range c.Rules {
		for _, macaddr := range t.Edge {
			edge := NewEdge(m.redis, macaddr)
			err = edge.SetUsed(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
