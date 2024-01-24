package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type edgeDeviceReq struct {
	CityHash string `json:"geo"`
	Offset   int    `json:"offset"`
	Num      int    `json:"num"`
}

type edgeDeviceRsp struct {
	Code int    `json:"code"`
	Meta string `json:"meta"`
	Data struct {
		Count int                  `json:"count"`
		Edges []edgeDeviceEdgeData `json:"edges"`
	} `json:"data"`
}

type edgeDeviceEdgeData struct {
	Macaddr string `json:"macaddr"`
	Public  string `json:"public"`
	ISP     string `json:"isp"`
	Single  int    `json:"single"`
}

func (m *implement) EdgeDevice(ctx context.Context, cityhash string, offset, num int) ([]Edge, error) {
	token, err := m.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(&edgeDeviceReq{
		CityHash: cityhash,
		Offset:   offset,
		Num:      num,
	})
	if err != nil {
		ctx.LogWithField("err", err).Error("json encode error")
		return nil, err
	}

	rsp, err := m.utils.Http().WithBody(reqBody).WithHeaders(map[string]string{"X-Token": token}).Post(fmt.Sprintf("%s/edge/device", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http post error")
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody := new(edgeDeviceRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return nil, err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("edge device error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("edge device error")
		return nil, err
	}

	data := make([]Edge, 0)
	for _, d := range rspBody.Data.Edges {
		edge := NewEdge(m.redis, d.Macaddr)
		data = append(data, edge)
		err = edge.SetDetail(ctx, &edgeDetail{
			Public:   d.Public,
			ISP:      d.ISP,
			Single:   d.Single,
			Cityhash: cityhash,
		})
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}
