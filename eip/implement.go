package eip

import (
	"github.com/picoorg/common/redis"
	"github.com/picoorg/common/utils"
)

func NewEIP(u utils.Utils, r redis.Redis, username, password, url, macaddr string) EIP {
	return &implement{
		utils:    u,
		redis:    r,
		username: username,
		password: password,
		url:      url,
		macaddr:  macaddr,
	}
}

type implement struct {
	utils    utils.Utils
	redis    redis.Redis
	username string
	password string
	url      string
	macaddr  string
}
