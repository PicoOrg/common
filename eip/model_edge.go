package eip

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/picoorg/common/context"
	"github.com/picoorg/common/redis"
)

type Edge interface {
	GetMacaddr() (m string)
	IsNew(ctx context.Context) (data bool, err error)
	SetNew(ctx context.Context) (err error)
	SetUsed(ctx context.Context) (err error)
	SetDetail(ctx context.Context, d *edgeDetail) (err error)
	GetDetail(ctx context.Context) (d *edgeDetail, err error)
}

func NewEdge(r redis.Redis, macaddr string) Edge {
	return &edgeModel{
		redis:   r,
		Macaddr: macaddr,
	}
}

func NewEdgeFromRedisKey(ctx context.Context, r redis.Redis, key string) (Edge, error) {
	value, exist, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	} else if !exist {
		err = fmt.Errorf("redis key %s not exist", key)
		ctx.LogWithField("err", err).Error("get redis eip edge error")
		return nil, err
	}
	return NewEdge(r, value), nil
}

type edgeModel struct {
	redis   redis.Redis
	Macaddr string `json:"macaddr"`
}

type edgeDetail struct {
	Public   string `json:"public"`
	ISP      string `json:"isp"`
	Single   int    `json:"single"`
	Cityhash string `json:"cityhash"`
}

func (m *edgeModel) GetMacaddr() string {
	return m.Macaddr
}

func (m *edgeModel) IsNew(ctx context.Context) (bool, error) {
	_, exist, err := m.redis.Get(ctx, m.getEdgeUsedRedisKey())
	if err != nil {
		return false, err
	} else if !exist {
		return true, nil
	}
	return false, nil
}

func (m *edgeModel) SetNew(ctx context.Context) error {
	err := m.redis.Del(ctx, m.getEdgeUsedRedisKey())
	if err != nil {
		return err
	}
	return m.redis.Set(ctx, m.getEdgeNewRedisKey(), m.Macaddr, 0)
}

func (m *edgeModel) SetUsed(ctx context.Context) error {
	err := m.redis.Del(ctx, m.getEdgeNewRedisKey())
	if err != nil {
		return err
	}
	return m.redis.Set(ctx, m.getEdgeUsedRedisKey(), m.Macaddr, 24*time.Hour)
}

func (m *edgeModel) GetDetail(ctx context.Context) (*edgeDetail, error) {
	key := m.getEdgeDetailRedisKey()
	data, exist, err := m.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	} else if !exist {
		err = fmt.Errorf("redis key %s not exist", key)
		ctx.LogWithField("err", err).Error("get redis eip edge error")
		return nil, err
	}
	d := new(edgeDetail)
	err = json.Unmarshal([]byte(data), &d)
	if err != nil {
		ctx.LogWithField("err", err).Error("json unmarshal error")
		return nil, err
	}
	return d, nil
}

func (m *edgeModel) SetDetail(ctx context.Context, d *edgeDetail) error {
	data, err := json.Marshal(d)
	if err != nil {
		ctx.LogWithField("err", err).Error("json marshal error")
		return err
	}
	return m.redis.Set(ctx, m.getEdgeDetailRedisKey(), string(data), 0)
}

func (m *edgeModel) getEdgeDetailRedisKey() string {
	return fmt.Sprintf("eip:edge:%s", m.Macaddr)
}

func (m *edgeModel) getEdgeNewRedisKey() string {
	return fmt.Sprintf("eip:edge:new:%s", m.Macaddr)
}

func (m *edgeModel) getEdgeUsedRedisKey() string {
	return fmt.Sprintf("eip:edge:used:%s", m.Macaddr)
}
