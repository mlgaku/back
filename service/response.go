package service

import (
	"encoding/json"
	"github.com/mlgaku/back/common"
	"github.com/mlgaku/back/types"
	"log"
)

type Response struct {
	Client *Client // 客户
}

// 打包数据
func (r *Response) Pack(pro types.Prot, val types.Value) []byte {
	pro.Body = common.StringValue(&val)
	b, _ := json.Marshal(&pro)
	return b
}

// 写内容
func (r *Response) Write(val []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("response failed: %s", r)
		}
	}()

	r.Client.Send <- val
}

// 获得 Response 实例
func NewResponse(cli *Client) *Response {
	return &Response{
		Client: cli,
	}
}
