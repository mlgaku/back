package main

import (
	"github.com/mlgaku/back/conf"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"log"
	"net/http"
)

func main() {

	// 路由
	APP.Route = conf.Route
	// 中间件
	APP.Middleware = conf.Middleware
	// 配置
	APP.Conf = NewConfig()
	// 发布订阅
	APP.Ps = NewPubsub()
	// 数据库
	APP.Db = NewDatabase()
	// Session
	APP.Session = NewSession(nil)

	// ws 服务
	APP.Server = NewServer()
	go APP.Server.Watch(handle)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		NewClient(APP.Server, w, r)
	})
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

// 处理消息
func handle(cli *Client, prim []byte) {
	mod, res := NewModule(cli), NewResponse(cli)

	val, err := mod.Load(prim)
	if err != nil {
		res.Write(res.Pack(*mod.Prot, &Fail{Msg: err.Error()}))
		return
	}
	if val != nil {
		res.Write(res.Pack(*mod.Prot, val))
	}
}
