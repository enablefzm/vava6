package dbstorage

import (
	"sync"
	"vava6/mysql"
	"vava6/valog"
	"vava6/vdb"
)

var dbStorage *vdb.DBStorage
var dbLk *sync.Mutex = new(sync.Mutex)

var CFG ConnCfg = ConnCfg{
	DBName:  "businessman",
	IPAdd:   "127.0.0.1",
	Port:    "3316",
	User:    "navuser01",
	Pass:    "fzmvava6",
	MaxConn: 6,
	MinConn: 2,
}

type ConnCfg struct {
	DBName  string
	IPAdd   string
	Port    string
	User    string
	Pass    string
	MaxConn int
	MinConn int
}

func GetStorage() (*vdb.DBStorage, error) {
	if dbStorage == nil {
		dbLk.Lock()
		defer dbLk.Unlock()
		if dbStorage == nil {
			// 创建MySql对象
			ptrMySql, err := mysql.NewDBs(
				CFG.DBName,
				CFG.IPAdd,
				CFG.Port,
				CFG.User,
				CFG.Pass,
				CFG.MaxConn,
				CFG.MinConn,
			)
			if err != nil {
				return nil, err
			}
			// 封装MySql
			pkMySql := vdb.NewMySql(ptrMySql)
			// 生成对象
			dbStorage = vdb.NewDBStorage(pkMySql)
			// LOG信息
			valog.OBLog.LogMessage("构造了dbStorage")
		}
	}
	return dbStorage, nil
}

func SaveOB(ob vdb.IFTargetDB) error {
	db, err := GetStorage()
	if err != nil {
		return err
	}
	return db.Save(ob)
}

func Querys(table string, keys map[string]interface{}) ([]map[string]string, error) {
	if conn, err := GetStorage(); err != nil {
		return nil, err
	} else {
		return conn.Querys(table, keys)
	}
}
