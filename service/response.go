package service

import (
	"github.com/mlgaku/back/common"
	"log"
)

type response struct {
	client *client
}

// 写内容
func (r *response) write(val common.Value) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("response failed: %s", r)
		}
	}()

	r.client.send <- common.BytesValue(&val)
}

// 创建替身
func (r *response) pseudo() *common.Response {
	return &common.Response{
		Write: func(val common.Value) {
			r.write(val)
		},
		Connection: r.client.connection,
	}
}

// 获得 response 实例
func newResponse(cli *client) *response {
	return &response{
		client: cli,
	}
}
