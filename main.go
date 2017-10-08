package main

import (
	"github.com/mlgaku/back/conf"
	. "github.com/mlgaku/back/service"
	"log"
	"net/http"
)

var app *App

func main() {
	app = &App{}

	// 路由
	app.Route = conf.Route
	// 配置
	app.Conf = NewConfig()
	// 发布订阅
	app.Ps = NewPubsub()
	// 数据库
	app.Db = NewDatabase(app.Conf.Db.Host, app.Conf.Db.Database, app.Conf.Db.Port)

	// ws 服务
	app.Server = NewServer()
	go app.Server.Watch(handle)

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		NewClient(app.Server, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 处理消息
func handle(cli *Client, prim []byte) {
	res, mod := NewResponse(cli), NewModule(app, cli)

	val, err := mod.Load(prim)
	if err != nil {
		res.Write([]byte(err.Error()))
	}

	res.Write(res.Pack(*mod.Prot, val))
}
