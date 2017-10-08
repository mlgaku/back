package service

var APP = &App{}

type App struct {
	Ps     *Pubsub
	Db     *Database
	Conf   *Config
	Server *Server

	Route map[string]interface{}
}
