package configuration

import (
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Server struct {
		Host string `toml:"ServerHost"`
		Port int    `toml:"ServerPort"`
	} `toml:"server"`
	Database struct {
		Host     string `toml:"DBHost"`
		Port     int    `toml:"DBPort"`
		User     string `toml:"DBUser"`
		Password string
		DBName   string `toml:"DBName"`
		DBTable  string `toml:"NewsTableName"`
		TTL      int    `toml:"TTL"`
	} `toml:"database"`
	News struct {
		BBCNewsRSSURL string        `toml:"BBCNewsRSSURL"`
		FetchInterval time.Duration `toml:"FetchInterval"`
		MaxRetries    int           `toml:"MaxRetries"`
	} `toml:"news"`
	Workers struct {
		NumWorkers int `toml:"NumWorkers"`
		ChanSize   int `toml:"ChanSize"`
	} `toml:"workers"`
}

var conf *Configuration
var once sync.Once

func NewConfiguration(configPath string) *Configuration {
	once.Do(func() {
		conf = &Configuration{}
		if _, err := toml.DecodeFile(configPath, conf); err != nil {
			panic(err)
		}
	})
	conf.Database.Password = os.Getenv("DB_PASSWORD")
	conf.News.FetchInterval = conf.News.FetchInterval * time.Minute
	return conf
}
