package main

import (
	"fmt"
	"io/ioutil"
	"log"
	mc "memcache"
	"route"
	"store"

	yaml "gopkg.in/yaml.v2"
)

var (
	config = GoBeansdbConfig{
		Listen:  "127.0.0.1",
		Port:    7900,
		WebPort: 7908,
		Threads: 4,
		ZK:      "",
	}
)

type GoBeansdbConfig struct {
	Addr    string `yaml:"-"` // HStoreConfig.Hostname:Port
	ZK      string `yaml:",omitempty"`
	Listen  string `yaml:",omitempty"`
	Port    int    `yaml:",omitempty"`
	WebPort int    `yaml:",omitempty"`
	Threads int    `yaml:",omitempty"`
	LogDir  string `yaml:",omitempty"`

	store.HStoreConfig `yaml:"hstore,omitempty"`
}

func loadServerConfig(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("read config failed", path, err.Error())
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		log.Fatal("unmarshal yaml format config failed")
	}
}

func loadConfigs(confdir string) {
	loadServerConfig(fmt.Sprintf("%s/%s", confdir, "server_global.yaml"))
	loadServerConfig(fmt.Sprintf("%s/%s", confdir, "server_local.yaml"))
	// route
	rt, err := route.LoadRouteTable(fmt.Sprintf("%s/%s", confdir, "route.yaml"), config.ZK)
	if err != nil {
		log.Fatalf("fail to load route table")
	}
	config.Addr = fmt.Sprintf("%s:%d", config.HStoreConfig.Hostname, config.Port)
	config.HStoreConfig.RouteConfig = rt.GetServerConfig(config.Addr)

	// config store
	config.HStoreConfig.InitForYaml()
	config.HStoreConfig.Init()
	store.SetConfig(config.HStoreConfig)
	// config mc
	mc.MaxKeyLength = config.HStoreConfig.MaxKeySize
	mc.MaxBodyLength = config.HStoreConfig.MaxValueSize
}