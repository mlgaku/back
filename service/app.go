package service

import (
	"log"
	"net/http"
)

// app 服务
type App struct{}

// 启动应用
func (*App) Start() {
	server := newServer()
	go server.watch()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		newClient(server, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
