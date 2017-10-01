package service

type request struct {
	body   []byte
	client *client
}

// 消息处理
func (r *request) handle() {
	res := newResponse(r.client)

	err := newModule(res).load(r.body)
	if err != nil {
		res.write(err.Error())
	}
}

// 获得 request 实例
func newRequest(cli *client, body []byte) *request {
	return &request{
		body:   body,
		client: cli,
	}
}
