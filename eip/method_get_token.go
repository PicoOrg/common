package eip

import (
	"fmt"
	"time"

	"github.com/picoorg/common/context"
)

func (m *implement) getTokenRedisKey() string {
	return fmt.Sprintf("eip:%s:token", m.macaddr)
}

func (m *implement) getConfigIDRedisKey() string {
	return fmt.Sprintf("eip:%s:config:id", m.macaddr)
}

func (m *implement) getConfigRulesRedisKey() string {
	return fmt.Sprintf("eip:%s:config:rules", m.macaddr)
}

func (m *implement) getEdgeNewAllRedisKey() string {
	return "eip:edge:new:*"
}

func (m *implement) GetToken(ctx context.Context) (string, error) {
	tokenRedisKey := m.getTokenRedisKey()
	token, exist, err := m.redis.Get(ctx, tokenRedisKey)
	if !exist {
		ctx.Info("no token found")
		token, err = m.clientAuth(ctx)
		if err != nil {
			return "", err
		}
		err = m.redis.Set(ctx, tokenRedisKey, token, 160*time.Hour)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return token, nil
}
