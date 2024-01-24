package eip

import (
	"fmt"

	"github.com/picoorg/common/context"
)

type RefreshTable struct {
	ID       int
	Cityhash string
}

func (m *implement) Refresh(ctx context.Context, tables []RefreshTable) (*gatewayConfig, error) {
	c, err := m.GatewayConfigGet(ctx, m.macaddr)
	if err != nil {
		return nil, err
	}

	err = m.CacheGatewayConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	rules := make([]gatewayConfigRule, 250)
	// for _, t := range c.Rules {
	// 	id, err := t.getIndex(ctx)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	rules[id-1] = t
	// }

	for _, t := range tables {
		var edges []Edge
		edges, err = m.getEdges(ctx, t.Cityhash, 4)
		if err != nil {
			return nil, err
		}
		rule := gatewayConfigRule{
			Enable:   true,
			Edge:     []string{},
			Network:  []string{fmt.Sprintf("172.30.168.%d", t.ID+1)},
			CityHash: t.Cityhash,
		}
		for _, e := range edges {
			// rule.Edge = append(rule.Edge, e.GetMacaddr())
			err = e.SetUsed(ctx)
			if err != nil {
				return nil, err
			}
		}
		rules[t.ID-1] = rule
	}
	c.ID = c.ID + 1
	c.Rules = make([]gatewayConfigRule, 0)
	for _, r := range rules {
		if r.Enable {
			c.Rules = append(c.Rules, r)
		}
	}

	_, err = m.gatewayConfigSet(ctx, m.macaddr, c)
	if err != nil {
		return nil, err
	}

	return c, m.CacheGatewayConfig(ctx, c)
}
