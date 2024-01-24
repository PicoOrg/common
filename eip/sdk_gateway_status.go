package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type gatewayStatusReq struct {
	Macaddr string `json:"macaddr"`
}

type gatewayStatusRsp struct {
	Code int             `json:"code"`
	Meta string          `json:"meta"`
	Data json.RawMessage `json:"data"`
}

func (m *implement) gatewayStatusDevice(ctx context.Context, macaddr string) (json.RawMessage, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(&gatewayStatusReq{Macaddr: macaddr})
	if err != nil {
		ctx.LogWithField("err", err).Error("json encode error")
		return nil, err
	}

	rsp, err := m.utils.Http().WithBody(reqBody).WithHeaders(map[string]string{"X-Token": token}).Post(fmt.Sprintf("%s/gateway/status", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http post error")
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody := new(gatewayStatusRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return nil, err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("gateway status error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("gateway status error")
		return nil, err
	}
	return rspBody.Data, nil
}
