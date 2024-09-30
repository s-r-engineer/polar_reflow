package configuration

type Config struct {
	Database Database `yaml:"database"`
	API      API      `yaml:"api"`
	Engine   Engine   `yaml:"engine"`
}

type Database struct {
	DBType   string `yaml:"host" envconfig:"DB_TYPE"`
	Host     string `yaml:"host" envconfig:"DB_HOST"`
	Token    string `yaml:"token" envconfig:"DB_TOKEN"`
	Database string `yaml:"org" envconfig:"DB_DATABASE"`
	Table    string `yaml:"bucket" envconfig:"DB_TABLE"`
	User     string `yaml:"user" envconfig:"DB_USER"`
	Password string `yaml:"password" envconfig:"DB_PASSWORD"`
}

type API struct {
	BindAddress string `yaml:"bindAddress" envconfig:"BIND_ADDRESS"`
}

type Engine struct {
	Parallel int `yaml:"parallel" envconfig:"PARALLELIZM"`
}
