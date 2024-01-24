package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type edgeStatusReq struct {
	Macaddr string `json:"macaddr"`
}

type edgeStatusRsp struct {
	Code int    `json:"code"`
	Meta string `json:"meta"`
	Data struct {
		Macaddr string `json:"macaddr"`
		Public  string `json:"public"`
		ISP     string `json:"isp"`
		Single  int    `json:"single"`
	} `json:"data"`
}

func (m *implement) edgeStatusDevice(ctx context.Context, macaddr string) (Edge, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(&edgeStatusReq{Macaddr: macaddr})
	if err != nil {
		ctx.LogWithField("err", err).Error("json encode error")
		return nil, err
	}

	rsp, err := m.utils.Http().WithBody(reqBody).WithHeaders(map[string]string{"X-Token": token}).Post(fmt.Sprintf("%s/edge/status", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http post error")
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody := new(edgeStatusRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return nil, err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("edge status error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("edge status error")
		return nil, err
	}

	edge := NewEdge(m.redis, macaddr)
	return edge, edge.SetDetail(ctx, &edgeDetail{
		Public:   rspBody.Data.Public,
		ISP:      rspBody.Data.ISP,
		Single:   rspBody.Data.Single,
		Cityhash: rspBody.Data.Macaddr,
	})
}
