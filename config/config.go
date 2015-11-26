package config

import (
	"fmt"
	"github.intra.douban.com/coresys/gobeansdb/utils"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

const (
	Version = "2.1.0.0"
)

var (
	DB    DBConfig
	Proxy ProxyConfig
	Route RouteTable
)

type ProxyConfig struct {
	ServerConfig `yaml:"server,omitempty"`
	MCConfig     `yaml:"mc,omitempty"`
}

type DBConfig struct {
	ServerConfig `yaml:"server,omitempty"`
	MCConfig     `yaml:"mc,omitempty"`
	HStoreConfig `yaml:"hstore,omitempty"`
}

type ServerConfig struct {
	Hostname string `yaml:",omitempty"`
	ZK       string `yaml:",omitempty"` // e.g. "zk1:2100"
	Listen   string `yaml:",omitempty"` // ip
	Port     int    `yaml:",omitempty"`
	WebPort  int    `yaml:",omitempty"`
	Threads  int    `yaml:",omitempty"` // NumCPU
	LogDir   string `yaml:",omitempty"`
}

func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

type MCConfig struct {
	MaxKeyLen int `yaml:"max_key_len,omitempty"`
	MaxReq    int `yaml:"max_req,omitempty"` // max num of requsets serve at the same time

	BodyMax int64 `yaml:"-"` // fail set/read_file if larger then this
	BodyBig int64 `yaml:"-"` // set may fail if memory is in shorage (determine by "storage")
	BodyInC int64 `yaml:"-"` // alloc body in cgo if larger then this

	BodyMaxStr string `yaml:"body_max_str,omitempty"`
	BodyBigStr string `yaml:"body_big_str,omitempty"`
	BodyInCStr string `yaml:"body_c_str,omitempty"`
}

func loadYamlConfig(config interface{}, path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, config)
}

func (c *DBConfig) checkEmptyConfig(path string) {
	if c.MaxKeyLen == 0 || c.SplitCapStr == "" || c.TreeHeight == 0 {
		log.Fatal("bad config: empty struct in ", path)
	}
}

func (c *DBConfig) Load(confdir string) {
	c.InitDefault()

	if confdir != "" {
		// global
		path := fmt.Sprintf("%s/%s", confdir, "beansdb_global.yaml")
		err := loadYamlConfig(c, path)
		if err != nil {
			log.Fatalf("bad config %s: %s", path, err.Error())
		}
		c.checkEmptyConfig(path)

		//local
		path = fmt.Sprintf("%s/%s", confdir, "beansdb_local.yaml")
		err = loadYamlConfig(c, path)
		if err != nil {
			log.Fatalf("bad config %s: %s", path, err.Error())
		}
		c.checkEmptyConfig(path)

		// route
		Route, err := LoadRouteTable(fmt.Sprintf("%s/%s", confdir, "route.yaml"), c.ZK)
		if err != nil {
			log.Fatalf("fail to load route table: %s", err.Error())
		}
		c.DBRouteConfig, err = Route.GetDBRouteConfig(c.Addr())
		if err != nil {
			log.Fatalf("bad config in %s %s", confdir, err.Error())
		}
	}
	utils.InitSizesPointer(c)
	err := c.HStoreConfig.init()
	if err != nil {
		log.Fatalf("bad config: %s", err.Error())
	}
}

func (c *DBConfig) InitDefault() {
	c.ServerConfig = DefaultServerConfig
	c.MCConfig = DefaultMCConfig
	c.HStoreConfig.InitDefault()
	utils.InitSizesPointer(c)
}

func (c *ProxyConfig) Load(confdir string) {
	// TODO:
}

func DumpConfig(config interface{}) {
	b, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("%s", err)
	} else {
		fmt.Println(string(b))
	}
}