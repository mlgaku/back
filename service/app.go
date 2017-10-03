package service

import (
	"fmt"
	"log"
	"net/http"
)

var (
	_db   *database
	_conf *config
)

// app 服务
type App struct{}

// 启动应用
func (*App) Start() {
	// 配置
	_conf = newConfig()

	// 数据库
	_db = newDatabase(fmt.Sprintf("%s:%d", _conf.conf.Db.Host, _conf.conf.Db.Port))
	defer _db.disconnect()

	// ws 服务
	server := newServer()
	go server.watch()
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		newClient(server, w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 获得 App 实例
func NewApp() *App {
	return &App{}
}
