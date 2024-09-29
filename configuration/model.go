package configuration

type Config struct {
	Database Database `yaml:"database"`
	API      API      `yaml:"api"`
	Engine   Engine   `yaml:"engine"`
}

type Mongo struct {
	Host       string `yaml:"host" envconfig:""`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

type Influx struct {
	Host   string `yaml:"host"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

type Database struct {
	Mongo  Mongo  `yaml:"mongo"`
	Influx Influx `yaml:"influx"`
}

type API struct {
	BindAddress string `yaml:"bindAddress"`
}

type Engine struct {
	Parallel int `yaml:"parallel"`
}
