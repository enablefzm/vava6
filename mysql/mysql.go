package mysql

import (
	"database/sql"
	"errors"
	"fmt"
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
	resDB.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&allowNativePasswords=true",
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
	// fmt.Println("构造了DB")
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

func (d *DBs) QuerySql(strSql string) ([]map[string]string, error) {
	rows, err := d.db.Query(strSql)
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

func (d *DBs) QuerysLimit(field, table, key string, page, limit int, orderby string) ([]map[string]string, error) {
	var sql string
	if page < 1 {
		page = 1
	}
	sql = fmt.Sprintf("SELECT %s FROM %s", field, table)
	if len(key) > 0 {
		sql = fmt.Sprintf("%s WHERE %s", sql, key)
	}
	if len(orderby) > 0 {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, orderby)
	}
	var idxLimit int = page*limit - limit
	sql = fmt.Sprintf("%s LIMIT %d, %d", sql, idxLimit, limit)
	// fmt.Println(sql)
	return d.QuerySql(sql)
}

// 查询多项数据库
func (d *DBs) Querys(field, table, key string) ([]map[string]string, error) {
	var sql string
	if len(key) > 2 {
		sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s", field, table, key)
	} else {
		sql = fmt.Sprintf("SELECT %s FROM %s LIMIT 0, 5000", field, table)
	}
	return d.QuerySql(sql)
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
			field = fmt.Sprint(field, ",", k, "=?")
		}
		i++
	}
	var where string
	j := 0
	for k, v := range key {
		if j > 0 {
			where += " AND "
		}
		switch value := v.(type) {
		case int, uint:
			where = fmt.Sprint(where, k, "=", value)
			//          where = fmt.Sprint(where, k, "=", strconv.Itoa(value))
			//		case uint:
			//			where = fmt.Sprint(where, k, "=", strconv.FormatUint(uint64(value), 10))
		case string:
			where = fmt.Sprint(where, k, "='", value, "'")
		default:
			return nil, errors.New("KEY参数错误")
		}
		j++
	}
	strSql := fmt.Sprint("UPDATE ", table, " SET ", field, " WHERE ", where) // "UPDATE " + table + " SET " + field + " WHERE " + where
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
