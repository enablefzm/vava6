package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL数据存储器
type DBs struct {
	dbName  string
	dbIP    string
	dbPort  string
	dbUser  string
	dbPass  string
	maxConn int
	minConn int
	db      *sql.DB
}

func NewDBs(dbName, dbIp, dbPort, dbUser, dbPass string, maxConn, minConn int) (*DBs, error) {
	resDB := &DBs{
		dbName:  dbName,
		dbIP:    dbIp,
		dbPort:  dbPort,
		dbUser:  dbUser,
		dbPass:  dbPass,
		maxConn: maxConn,
		minConn: minConn,
	}
	var err error
	resDB.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		resDB.dbUser,
		resDB.dbPass,
		resDB.dbIP,
		resDB.dbPort,
		resDB.dbName,
	))
	if err != nil {
		return resDB, err
	}
	if err = resDB.db.Ping(); err != nil {
		return resDB, err
	}
	resDB.db.SetMaxOpenConns(maxConn)
	resDB.db.SetMaxIdleConns(minConn)
	fmt.Println("构造了DB")
	return resDB, err
}

// 查询数据库信息
func (d *DBs) NaviQuery(field, table, key string) (map[string]string, error) {
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", field, table, key)
	// fmt.Println("SQLDB Query:", sql)
	rows, err := d.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			} else {
				record[columns[i]] = ""
			}
		}
	}
	return record, nil
}

// 查询多项数据库
func (d *DBs) Querys(field, table, key string) ([]map[string]string, error) {
	var sql string
	if len(key) > 2 {
		sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s", field, table, key)
	} else {
		sql = fmt.Sprintf("SELECT %s FROM %s LIMIT 0, 5000", field, table)
	}
	rows, err := d.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	records := make([]map[string]string, 0, 10)
	for rows.Next() {
		record := make(map[string]string)
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			} else {
				record[columns[i]] = ""
			}
		}
		records = append(records, record)
	}
	return records, nil
}

// 插入数据到数据库
func (d *DBs) Insert(table string, info map[string]interface{}) (sql.Result, error) {
	field := ""
	value := ""
	vArr := make([]interface{}, len(info))
	i := 0
	for k, v := range info {
		// kArr[i] = k
		vArr[i] = v
		if len(field) < 1 {
			field = k
			value = "?"
		} else {
			field += "," + k
			value += ",?"
		}
		i++
	}
	sql := "INSERT INTO " + table + "(" + field + ") VALUES(" + value + ")"
	// fmt.Println("SQLDB Insert:", sql)
	stmt, err := d.db.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := stmt.Exec(vArr...)
	return res, err
}

// 更新数据到数据库
func (d *DBs) Update(table string, info, key map[string]interface{}) (sql.Result, error) {
	il := len(info)
	if il < 1 || len(key) < 1 {
		return nil, errors.New("UPDATE参数错误")
	}
	i := 0
	vArr := make([]interface{}, il)
	field := ""
	for k, v := range info {
		vArr[i] = v
		if len(field) < 1 {
			field = k + "=?"
		} else {
			field += "," + k + "=?"
		}
		i++
	}
	var where string
	j := 0
	for k, v := range key {
		value, ok := v.(int)
		if j > 0 {
			where += " AND "
		}
		if ok {
			where += k + "=" + strconv.Itoa(value)
		} else if sv, ok := v.(string); ok {
			where += k + "='" + sv + "'"
		} else {
			return nil, errors.New("KEY参数错误")
		}
		j++
		// vArr[i] = v
	}
	strSql := "UPDATE " + table + " SET " + field + " WHERE " + where
	// fmt.Println("SQLDB Update:", strSql)
	stmt, err := d.db.Prepare(strSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(vArr...)
	return res, err
}

// 删除数据
func (d *DBs) Remove(table string, key map[string]interface{}) (sql.Result, error) {
	il := len(key)
	if il < 1 {
		return nil, errors.New("没有删除条件")
	}
	i := 0
	vArr := make([]interface{}, il)
	arrKey := make([]string, 0, il)
	for k, v := range key {
		vArr[i] = v
		arrKey = append(arrKey, fmt.Sprint(k, "=?"))
		i++
	}
	strSql := fmt.Sprint("DELETE FROM ", table, " WHERE ", strings.Join(arrKey, " AND "))
	fmt.Println(strSql)
	stmt, err := d.db.Prepare(strSql)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(vArr...)
	return res, err

}

func (d *DBs) GetDB() *sql.DB {
	return d.db
}
