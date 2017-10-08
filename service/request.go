package service

type Request struct {
	Body   []byte  // 内容
	Client *Client // 客户
}

// 获得 Request 实例
func NewRequest(body []byte, cli *Client) *Request {
	return &Request{
		Body:   body,
		Client: cli,
	}
}
