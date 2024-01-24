package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type edgeCityReq struct {
}

type edgeCityRsp struct {
	Code int                          `json:"code"`
	Meta string                       `json:"meta"`
	Data map[string]map[string]string `json:"data"`
}

func (m *implement) edgeCity(ctx context.Context) (map[string]map[string]string, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := m.utils.Http().WithHeaders(map[string]string{"X-Token": token}).Get(fmt.Sprintf("%s/edge/city", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http get error")
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody := new(edgeCityRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return nil, err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("edge city error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("edge city error")
		return nil, err
	}
	return rspBody.Data, nil
}
