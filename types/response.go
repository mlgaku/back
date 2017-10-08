package types

type Response struct {
	Write  func([]byte) // 写内容
	Client *Client      // 客户
}
