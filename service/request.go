package service

import "github.com/mlgaku/back/types"

type request struct {
	body   []byte
	client *client
}

// 消息处理
func (r *request) handle() {
	res := newResponse(r.client)

	err := newModule(r, res).load(r.body)
	if err != nil {
		res.write(err.Error())
	}
}

// 创建替身
func (r *request) pseudo() *types.Request {
	return &types.Request{
		Body:    r.body,
		BodyStr: string(r.body),
	}
}

// 获得 request 实例
func newRequest(cli *client, body []byte) *request {
	return &request{
		body:   body,
		client: cli,
	}
}
