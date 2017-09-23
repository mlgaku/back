package service

type request struct{}

// 消息处理
func (*request) handle(cli *client, msg []byte) {
	new(module).init(new(response).init(cli)).load(msg)
}
