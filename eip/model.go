package eip

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/picoorg/common/context"
)

type gateway struct {
	Macaddr     string `json:"macaddr"`
	Customer    string `json:"customer"`
	Information string `json:"information"`
	Enable      bool   `json:"enable"`
}

type gatewayConfig struct {
	ID    int                 `json:"id"`
	Rules []gatewayConfigRule `json:"rules"`
}

type gatewayConfigRule struct {
	Table    int      `json:"table"`
	Enable   bool     `json:"enable"`
	Edge     []string `json:"edge"`
	Network  []string `json:"network"`
	CityHash string   `json:"cityhash"`
}

func (r *gatewayConfigRule) getIndex(ctx context.Context) (id int, err error) {
	if len(r.Network) == 0 {
		err = fmt.Errorf("no network")
		ctx.LogWithField("err", err).Error("no network")
		return
	}
	ipv4 := strings.Split(r.Network[0], ".")
	if len(ipv4) != 4 {
		err = fmt.Errorf("no ipv4")
		ctx.LogWithField("err", err).Error("no ipv4")
		return
	}
	var id64 int64
	id64, err = strconv.ParseInt(ipv4[3], 10, 32)
	if err != nil {
		ctx.LogWithField("err", err).Error("parse int error")
		return
	}
	id = int(id64) - 1
	return
}
