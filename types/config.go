package types

type Config struct {
	Db struct {
		Host     string
		Port     int
		Database string
	}
	App struct {
		Debug bool
	}
	Secret struct {
		Salt string
	}
}
