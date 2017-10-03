package types

type Db struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

type App struct {
	Debug bool `json:"debug"`
}

type Config struct {
	Db  Db  `json:"db"`
	App App `json:"app"`
}
