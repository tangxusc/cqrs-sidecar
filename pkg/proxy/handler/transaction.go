package handler

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"github.com/tangxusc/cqrs-sidecar/pkg/proxy"
	"strings"
)

func init() {
	proxy.Handlers = append(proxy.Handlers, &transaction{})
}

type transaction struct {
}

func (s *transaction) Match(sql string) bool {
	space := strings.TrimSpace(sql)
	space = strings.ToLower(space)
	if space == `begin` || space == `commit` || space == `rollback` || space == `start transaction` {
		return true
	}

	return false
}

func (s *transaction) Handler(query string, handler *proxy.ConnHandler) (*mysql.Result, error) {
	query = strings.TrimSpace(query)
	query = strings.ToLower(query)
	switch query {
	case `begin`, `start transaction`:
		tx, e := db.ConnInstance.Begin()
		if e != nil {
			return nil, e
		}
		handler.Tx = tx
	case `commit`:
		if handler.Tx == nil {
			return nil, fmt.Errorf(`需要先开启事务,才能提交事务`)
		}
		e := handler.Tx.Commit()
		return nil, e
	case `rollback`:
		if handler.Tx == nil {
			return nil, fmt.Errorf(`需要先开启事务,才能回滚事务`)
		}
		e := handler.Tx.Rollback()
		return nil, e
	}
	return nil, nil
}
