package main

import (
	"github.com/mlgaku/back/conf"
	. "github.com/mlgaku/back/service"
	"log"
	"net/http"
)

func main() {

	// 路由
	APP.Route = conf.Route
	// 配置
	APP.Conf = NewConfig()
	// 发布订阅
	APP.Ps = NewPubsub()
	// 数据库
	APP.Db = NewDatabase()

	// ws 服务
	APP.Server = NewServer()
	go APP.Server.Watch(handle)

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		NewClient(APP.Server, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 处理消息
func handle(cli *Client, prim []byte) {
	res, mod := NewResponse(cli), NewModule(cli)

	val, err := mod.Load(prim)
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	if val != nil {
		res.Write(res.Pack(*mod.Prot, val))
	}
}
