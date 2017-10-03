package types

// 路由信息
type Route struct {
	Module string `json:"mod"`  // 模块
	Action string `json:"act"`  // 操作
	Body   string `json:"body"` // 正文
}
