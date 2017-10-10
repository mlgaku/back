package service

var APP = &App{}

type App struct {
	Ps      *Pubsub
	Db      *Database
	Conf    *Config
	Server  *Server
	Session *Session

	Route      map[string]interface{}
	Middleware map[string][]interface{}
}
