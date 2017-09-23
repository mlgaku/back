package service

import (
	"github.com/sxyazi/maile/common"
	"log"
)

type response struct {
	client *client
}

// 初始化
func (r *response) init(c *client) *response {
	r.client = c
	return r
}

// 写内容
func (r *response) write(v common.Value) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("response failed: %s", r)
		}
	}()

	r.client.send <- common.BytesValue(&v)
}
