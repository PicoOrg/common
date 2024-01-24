package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type gatewayConfigSetReq struct {
	Macaddr string         `json:"macaddr"`
	Config  *gatewayConfig `json:"config"`
}

type gatewayConfigSetRsp struct {
	Code int    `json:"code"`
	Meta string `json:"meta"`
	Data string `json:"data"`
}

func (m *implement) gatewayConfigSet(ctx context.Context, macaddr string, config *gatewayConfig) (string, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return "", err
	}

	reqBody, err := json.Marshal(&gatewayConfigSetReq{Macaddr: macaddr, Config: config})
	if err != nil {
		ctx.LogWithField("err", err).Error("json encode error")
		return "", err
	}

	rsp, err := m.utils.Http().WithBody(reqBody).WithHeaders(map[string]string{"X-Token": token}).Post(fmt.Sprintf("%s/gateway/config/set", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http post error")
		return "", err
	}
	defer rsp.Body.Close()

	rspBody := new(gatewayConfigSetRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return "", err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("gateway config set error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("gateway config set error")
		return "", err
	}
	return rspBody.Data, nil
}
