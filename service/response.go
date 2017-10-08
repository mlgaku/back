package service

import (
	"github.com/mlgaku/back/types"
	"log"
)

type response struct {
	client *client // 客户
}

// 写内容
func (r *response) write(val []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("response failed: %s", r)
		}
	}()

	r.client.send <- val
}

// 创建替身
func (r *response) pseudo() *types.Response {
	return &types.Response{
		Write: func(v []byte) {
			r.write(v)
		},
		Client: r.client.pseudo(),
	}
}

// 获得 response 实例
func newResponse(cli *client) *response {
	return &response{
		client: cli,
	}
}
