package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type gatewayConfigGetReq struct {
	Macaddr string `json:"macaddr"`
}

type gatewayConfigGetRsp struct {
	Code    int    `json:"code"`
	Meta    string `json:"meta"`
	RawData []byte `json:"data"`
}

func (m *implement) GatewayConfigGet(ctx context.Context, macaddr string) (*gatewayConfig, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(&gatewayConfigGetReq{Macaddr: macaddr})
	if err != nil {
		ctx.LogWithField("err", err).Error("json encode error")
		return nil, err
	}

	rsp, err := m.utils.Http().WithBody(reqBody).WithHeaders(map[string]string{"X-Token": token}).Post(fmt.Sprintf("%s/gateway/config/get", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http post error")
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody := new(gatewayConfigGetRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return nil, err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("gateway config get error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("gateway config get error")
		return nil, err
	} else if len(rspBody.RawData) == 0 || rspBody.RawData == nil {
		rspBody.RawData = []byte("{}")
	}

	data := new(gatewayConfig)
	err = json.Unmarshal(rspBody.RawData, &data)
	if err != nil {
		ctx.LogWithField("err", err).Error("json unmarshal error")
		return nil, err
	}
	return data, nil
}
