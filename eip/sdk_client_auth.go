package eip

import (
	"encoding/json"
	"fmt"

	"github.com/picoorg/common/context"
)

type clientAuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type clientAuthRsp struct {
	Code int    `json:"code"`
	Meta string `json:"meta"`
	Data string `json:"data"`
}

func (m *implement) clientAuth(ctx context.Context) (string, error) {
	reqBody, err := json.Marshal(&clientAuthReq{
		Username: m.username,
		Password: m.password,
	})
	if err != nil {
		ctx.LogWithField("err", err).Error("json encode error")
		return "", err
	}

	rsp, err := m.utils.Http().WithBody(reqBody).Post(fmt.Sprintf("%s/client/auth", m.url))
	if err != nil {
		ctx.LogWithField("err", err).Error("http post error")
		return "", err
	}
	defer rsp.Body.Close()

	rspBody := new(clientAuthRsp)
	err = json.NewDecoder(rsp.Body).Decode(rspBody)
	if err != nil {
		ctx.LogWithField("err", err).Error("json decode error")
		return "", err
	}

	if rspBody.Code != 0 {
		err = fmt.Errorf("client auth error: %s", rspBody.Meta)
		ctx.LogWithField("err", err).Error("client auth error")
		return "", err
	}
	return rspBody.Data, nil
}
