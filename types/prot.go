package types

// 交换协议
type Prot struct {
	Mod  string `json:"mod"`  // 模块
	Act  string `json:"act"`  // 行为
	Body string `json:"body"` // 正文
}
