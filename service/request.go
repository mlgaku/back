package service

type Request struct {
	Body   []byte  // 内容
	Client *Client // 客户
}

// 获取远端地址
func (r *Request) RemoteAddr() string {
	if ip := r.Client.Http.Header.Get("X-Real-IP"); ip != "" {
		return ip
	} else if ip := r.Client.Http.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	} else {
		return r.Client.Http.RemoteAddr
	}
}

// 获得 Request 实例
func NewRequest(body []byte, cli *Client) *Request {
	return &Request{
		Body:   body,
		Client: cli,
	}
}
