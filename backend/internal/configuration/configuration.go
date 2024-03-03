package configuration

import (
	"github.com/BurntSushi/toml"
	"sync"
)

type Configuration struct {
	Server struct {
		Port int `toml:"ServerPort"`
	} `toml:"server"`
	News struct {
		BBCNewsRSSURL string `toml:"BBCNewsRSSURL"`
	} `toml:"news"`
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
	return conf
}
