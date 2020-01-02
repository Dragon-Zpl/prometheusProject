package conf

import (
	"github.com/astaxie/beego/logs"
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"strings"
)

var (
	PromeHotReloadCfg HotReload
	PromethuesPath  PromePath
	Email  EmailConfig
	Consul ConsulConfig
)

type HotReload struct {
	PromeUrl string
	AlertUrl string
}

type PromePath struct {
	Path string
	Recieve string
}

type EmailConfig struct {
	Sender string
	Host string
	Port int
	Password string
}

type ConsulConfig struct {
	Url string
	Prefix string
	DirName string
}

func GetRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Error(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func InitConf() {
	confPath := GetRootPath() + "/conf/conf.ini"
	cfg, err := ini.Load(confPath)
	if err != nil {
		panic(err)
	}

	err = cfg.Section("PromeHotRelood").MapTo(&PromeHotReloadCfg)
	if err != nil {
		logs.Error("cfg.MapTo PromeHotRelood settings err: %v", err)
	}

	err = cfg.Section("PromethuesPath").MapTo(&PromethuesPath)
	if err != nil {
		logs.Error("cfg.MapTo PromethuesPath settings err: %v", err)
	}


	err = cfg.Section("Email").MapTo(&Email)
	if err != nil {
		logs.Error("cfg.MapTo Email settings err: %v", err)
	}

	err = cfg.Section("Consul").MapTo(&Consul)
	if err != nil {
		logs.Error("cfg.MapTo Consul settings err: %v", err)
	}
}
