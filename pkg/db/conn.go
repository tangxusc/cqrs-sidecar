package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-sidecar/pkg/config"
	"os"
	"time"
)

type Conn struct {
	*sql.DB
}

var ConnInstance *Conn

func InitConn(ctx context.Context) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=true", config.Instance.Db.Username, config.Instance.Db.Password,
		"tcp", config.Instance.Db.Address, config.Instance.Db.Port, config.Instance.Db.Database)
	var e error
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		logrus.Errorf("[proxy]连接出现错误,url:%v,错误:%v", dsn, e.Error())
		os.Exit(1)
	}
	db.SetConnMaxLifetime(time.Duration(config.Instance.Db.LifeTime) * time.Second)
	db.SetMaxOpenConns(config.Instance.Db.MaxOpen)
	db.SetMaxIdleConns(config.Instance.Db.MaxIdle)
	ConnInstance = &Conn{db}
}

func CloseConn() {
	if ConnInstance != nil {
		ConnInstance.Close()
		ConnInstance = nil
	}
}

func (conn *Conn) ExecWithTx(tx *sql.Tx, sqlString string, param ...interface{}) error {
	logrus.Debugf("[proxy]Exec:%s,param:%v", sqlString, param)
	stmt, e := tx.Prepare(sqlString)
	if e != nil {
		return e
	}
	defer stmt.Close()
	_, e = stmt.Exec(param...)
	return e
}

func (conn *Conn) Exec(sqlString string, param ...interface{}) error {
	logrus.Debugf("[proxy]Exec:%s,param:%v", sqlString, param)
	return conn.Tx(func(tx *sql.Tx) error {
		stmt, e := tx.Prepare(sqlString)
		if e != nil {
			return e
		}
		defer stmt.Close()
		_, e = stmt.Exec(param...)
		return e
	})
}

func (conn *Conn) Execs(sqlString string, params [][]interface{}) error {
	logrus.Debugf("[proxy]Execs:%s,params:%v", sqlString, params)
	return conn.Tx(func(tx *sql.Tx) error {
		stmt, e := tx.Prepare(sqlString)
		if e != nil {
			return e
		}
		defer stmt.Close()
		for _, value := range params {
			_, e := stmt.Exec(value...)
			if e != nil {
				return e
			}
		}
		return nil
	})
}

func (conn *Conn) Tx(f func(tx *sql.Tx) error) error {
	logrus.Debugf("[proxy]Tx:%v", f)
	tx, e := conn.Begin()
	if e != nil {
		return e
	}
	e = f(tx)
	if e != nil {
		defer tx.Rollback()
		return e
	}
	return tx.Commit()
}

func (conn *Conn) QueryOne(sqlString string, scan []interface{}, param ...interface{}) error {
	logrus.Debugf("[proxy]QueryOne:%s,param:%v", sqlString, param)
	stmt, e := conn.Prepare(sqlString)
	if e != nil {
		return e
	}
	defer stmt.Close()
	row := stmt.QueryRow(param...)
	e = row.Scan(scan...)
	//未找到记录
	if e != nil && e == sql.ErrNoRows {
		return nil
	}
	if e != nil {
		return e
	}
	return nil
}

func (conn *Conn) QueryList(sqlString string, newRow func(types []*sql.ColumnType) []interface{}, param ...interface{}) error {
	return conn.Query(sqlString, newRow, func(row []interface{}) {
	}, func(strings []string) {
	}, param...)
}

func (conn *Conn) Query(query string, newRow func(types []*sql.ColumnType) []interface{}, rowAfter func(row []interface{}), setColumnNames func([]string), param ...interface{}) error {
	logrus.Debugf("[proxy]Query:%s,param:%v", query, param)
	stmt, e := conn.Prepare(query)
	if e != nil {
		return e
	}
	rows, e := stmt.Query(param...)
	if e != nil {
		return e
	}
	defer rows.Close()
	types, e := rows.ColumnTypes()
	if e != nil {
		return e
	}
	strings, e := rows.Columns()
	if e != nil {
		return e
	}
	setColumnNames(strings)
	for rows.Next() {
		row := newRow(types)
		e := rows.Scan(row...)
		if e != nil {
			return e
		}
		rowAfter(row)
	}
	return nil
}

func (conn *Conn) Proxy(query string) (columnNames []string, columnValues [][]interface{}, err error) {
	logrus.Debugf("[proxy]Proxy:%s", query)
	var temp interface{} = ""
	var rowOrigin []interface{}
	var result []interface{}

	columnValues = make([][]interface{}, 0)
	err = conn.Query(query,
		func(types []*sql.ColumnType) []interface{} {
			if result == nil {
				result = make([]interface{}, len(types))
				rowOrigin = make([]interface{}, 0, len(types))
				for key := range types {
					rowOrigin = append(rowOrigin, temp)
					result[key] = &rowOrigin[key]
				}
			}
			return result
		},
		func(row []interface{}) {
			i := make([]interface{}, len(row))
			for key := range row {
				v1 := rowOrigin[key]
				switch v1.(type) {
				case time.Time:
					i[key] = v1.(time.Time).String()
					continue
				}
				i[key] = rowOrigin[key]
			}
			columnValues = append(columnValues, i)
		},
		func(strings []string) {
			columnNames = strings
		})
	if err != nil {
		return
	}
	return
}

/*
保存事件
根据id查询,如果存在,则ack
如果不存在则插入,然后ack
*/
func (conn *Conn) Save(Id string, EventType string, AggId string, AggType string, CreateTime time.Time, Data string, status string) error {
	tx, e := ConnInstance.DB.Begin()
	if e != nil {
		return e
	}
	result := tx.QueryRow(`select count(1) c from event where id =? `, Id)
	var count int
	e = result.Scan(&count)
	if e != nil {
		_ = tx.Commit()
		return e
	}
	//已存在event
	if count > 0 {
		_ = tx.Commit()
		return nil
	}
	e = ConnInstance.ExecWithTx(tx, `insert into event(id,event_type,agg_id,agg_type,create_time,data,status) values(?,?,?,?,?,?,?)`,
		Id, EventType, AggId, AggType, CreateTime, Data, status)
	if e != nil {
		r := tx.Rollback()
		if r != nil {
			return fmt.Errorf(`%v,%v`, e, r)
		} else {
			return e
		}
	}
	e = tx.Commit()
	return e
}
