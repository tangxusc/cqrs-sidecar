package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"github.com/tangxusc/cqrs-sidecar/pkg/proxy"
)

type defaultHandler struct {
}

func (d *defaultHandler) Match(stmt string) bool {
	return true
}

func (d *defaultHandler) Handler(query string, handler *proxy.ConnHandler) (*mysql.Result, error) {
	columnNames, columnValues, err := db.ConnInstance.Proxy(query)
	if err != nil {
		return nil, err
	}
	resultSet, err := mysql.BuildSimpleTextResultset(columnNames, columnValues)
	if err != nil {
		return nil, err
	}

	return &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultSet,
	}, err
}

func init() {
	proxy.DefaultHandler = &defaultHandler{}
}
