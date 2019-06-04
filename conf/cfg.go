package conf

import (
	"encoding/json"
	"eyeSky/modules/transport/util"
	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	VERSION = "0.0.1"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type NSQConfig struct {
	Enabled      bool     `json:"enabled"`
	Addrs        []string `json:"addrs"`
	DialTimeout  int      `json:"dial_timeout"`
	ReadTimeout  int      `json:"read_timeout"`
	WriteTimeout int      `json:"write_timeout"`
	ClientId     string   `json:"client_id"`
	UserAgent    string
	AuthSecret   string
	Retry        int
	//DefaultTopic string
	Topics     []string
	EventTopic string
	Batch      int
}

func NsqConfig() *nsq.Config {
	return nsq.NewConfig()
}

type GlobalConfig struct {
	LogLevel int        `json:"loglevel"`
	Http     HttpConfig `json:"http"`
	Rpc      RpcConfig  `json:"rpc"`
	NSQ      NSQConfig  `json:"nsq"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()

	return config
}

func LoadConfig(cfgPath string) {
	if cfgPath == "" {
		log.Fatal("cfgPath is null")
	}

	if !util.IsExits(cfgPath) {
		log.Fatalf("%s is not existed", cfgPath)
	}

	configContent, err := util.ToTrimString(cfgPath)
	if err != nil {
		log.Fatal("read config file:", cfgPath, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatal("parse config file:", cfgPath, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	ConfigFile = cfgPath

	log.Infoln("conf.ParseConfig ok, file", cfgPath)

	if config.NSQ.Batch <= 0 {
		config.NSQ.Batch = 500
	}

	if config.NSQ.Retry <= 0 {
		config.NSQ.Retry = 3
	}
}
