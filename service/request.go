package service

import "github.com/mlgaku/back/types"

type request struct {
	body   string  // 内容
	prim   []byte  // 原始内容
	client *client // 客户
}

// 消息处理
func (r *request) handle() {
	res := newResponse(r.client)

	err := newModule(r, res).load(r.prim)
	if err != nil {
		res.write([]byte(err.Error()))
	}
}

// 创建替身
func (r *request) pseudo() *types.Request {
	return &types.Request{
		Body:    []byte(r.body),
		BodyStr: r.body,
	}
}

// 获得 request 实例
func newRequest(cli *client, prim []byte) *request {
	return &request{
		prim:   prim,
		client: cli,
	}
}
