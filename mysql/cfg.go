package mysql

import (
	"vava6/vaini"
	"vava6/vatools"
)

// 获取配置信息
type Cfg struct {
	User    string
	Pass    string
	Address string
	Port    string
	DBName  string
	MaxConn int
	MinConn int
}

func NewCfg() *Cfg {
	cfg := &Cfg{
		User:    "navuser01",
		Pass:    "fzmvava6",
		Address: "127.0.0.1",
		Port:    "3316",
		DBName:  "twic",
		MaxConn: 10,
		MinConn: 3,
	}
	// 加载文件
	path, err := vatools.GetNowPath()
	if err == nil {
		c := vaini.NewConfig(path + "/cfg.ini")
		if mp, ok := c.GetNode("database"); ok {
			for k, v := range mp {
				switch k {
				case "user":
					cfg.User = v
				case "pass":
					cfg.Pass = v
				case "address":
					cfg.Address = v
				case "port":
					cfg.Port = v
				case "dbname":
					cfg.DBName = v
				case "maxconn":
					cfg.MaxConn = vatools.SInt(v)
				case "minconn":
					cfg.MinConn = vatools.SInt(v)
				}
			}
		}
	}
	return cfg
}
