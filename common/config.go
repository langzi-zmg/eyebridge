package common

import (
	std "gitlab.wallstcn.com/wscnbackend/ivankastd"
	"gitlab.wallstcn.com/operation/eyebridge/models"
)

var (
	GlobalConf *Config
)

type Config struct {
	Micro     std.ConfigService `yaml:"micro"`
	Log       std.ConfigLog     `yaml:"log"`
	Mysql     std.ConfigMysql   `yaml:"mysql"`
	EtcdAddrs EtcdAddrsDetail   `yaml:"etcd_addrs"`
	Redis     Redis             `yaml:"redis"`
	Bind      string            `yaml:"bind"`
	CertPem   string            `yaml:"certpem"`
	KeyPem    string            `yaml:"keypem"`
}
type EtcdAddrsDetail struct {
	SH []string `yaml:"sh"`
	BJ []string `yaml:"bj"`
	GZ []string `yaml:"gz"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password""`
	DB       string `yaml:"db"`
}

func LoadConfig(filePath string) {
	println("loading config")
	GlobalConf = &Config{}
	std.LoadConf(GlobalConf, filePath)
}

func Initalise() {
	std.InitLog(GlobalConf.Log)
	models.InitModel(GlobalConf.Mysql)
}
