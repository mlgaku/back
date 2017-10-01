package service

type request struct{}

// 消息处理
func (*request) handle(cli *client, msg []byte) {
	res := new(response).init(cli)

	err := new(module).init(res).load(msg)
	if err != nil {
		res.write(err.Error())
	}
}
