package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type gatewayListReq struct {
}

type gatewayListRsp struct {
	Code int       `json:"code"`
	Meta string    `json:"meta"`
	Data []gateway `json:"data"`
}

func (m *implement) GatewayList(ctx context.Context) ([]gateway, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := m.utils.Http().WithHeaders(map[string]string{"X-Token": token}).Get(fmt.Sprintf("%s/gateway/list", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http get error")
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody := new(gatewayListRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return nil, err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("gateway list error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("gateway list error")
		return nil, err
	}
	return rspBody.Data, nil
}
