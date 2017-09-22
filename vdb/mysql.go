package vdb

// 将MySql封装成IFStorage接口对象

import (
	"fmt"
	"strconv"
	"strings"
	"vava6/mysql"
)

func NewMySql(pt *mysql.DBs) *MySql {
	return &MySql{
		ptrDBs: pt,
	}
}

type MySql struct {
	ptrDBs *mysql.DBs
}

func (this *MySql) Querys(table string, keys map[string]interface{}) ([]map[string]string, error) {
	arrKey := make([]string, 0, len(keys))
	for k, v := range keys {
		if value, ok := v.(int); ok {
			arrKey = append(arrKey, k+"="+strconv.Itoa(value))
		} else if sv, ok := v.(string); ok {
			arrKey = append(arrKey, k+"='"+sv+"'")
		} else {
			return nil, fmt.Errorf("querys key error.")
		}
	}
	strKey := strings.Join(arrKey, " AND ")
	return this.ptrDBs.Querys("*", table, strKey)
}

func (this *MySql) Insert(table string, info map[string]interface{}) (int64, error) {
	if res, err := this.ptrDBs.Insert(table, info); err != nil {
		return 0, err
	} else {
		if idx, err := res.LastInsertId(); err != nil {
			return 0, nil
		} else {
			return idx, nil
		}
	}
}

func (this *MySql) Update(table string, info, keys map[string]interface{}) (int64, error) {
	res, err := this.ptrDBs.Update(table, info, keys)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (this *MySql) Remove(table string, keys map[string]interface{}) (int64, error) {
	res, err := this.ptrDBs.Remove(table, keys)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
