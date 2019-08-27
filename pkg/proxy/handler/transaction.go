package handler

import (
	"github.com/siddontang/go-mysql/mysql"
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
	if space == `begin` || space == `commit` || space == `rollback` {
		return true
	}

	return false
}

func (s *transaction) Handler(query string, handler *proxy.ConnHandler) (*mysql.Result, error) {
	//TODO:事务处理
	//开启事务 handler.tx=tx,开启事务后,ack直接使用这个tx
	//回滚事务 handler.tx.rollback()
	//提交事务 handler.tx.commit()
	return nil, nil
}
