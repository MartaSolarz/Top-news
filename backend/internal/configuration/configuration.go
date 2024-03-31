package configuration

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Server struct {
		Host     string `toml:"ServerHost"`
		Port     int    `toml:"ServerPort"`
		MainPath string `toml:"MainPath"`
	} `toml:"server"`
	Database struct {
		Host       string `toml:"DBHost"`
		Port       int    `toml:"DBPort"`
		User       string `toml:"DBUser"`
		Password   string
		DBName     string `toml:"DBName"`
		NewsTable  string `toml:"NewsTableName"`
		EmailTable string `toml:"EmailTableName"`
		TTL        int    `toml:"TTL"`
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
	Email struct {
		HostEmail  string `toml:"HostEmail"`
		HostPass   string
		EmailTopic string `toml:"EmailTopic"`
	}
	OpenAPI struct {
		Url     string `toml:"APIUrl"`
		Key     string
		Disable bool
	}
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
	conf.Email.HostPass = os.Getenv("EMAIL_PASSWORD")
	conf.OpenAPI.Key = os.Getenv("OPEN_API_KEY")

	disable, err := strconv.ParseBool(os.Getenv("DISABLE_OPEN_API"))
	if err != nil {
		disable = true
	}
	conf.OpenAPI.Disable = disable

	return conf
}
