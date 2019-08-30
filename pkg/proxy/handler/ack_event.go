package handler

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"github.com/tangxusc/cqrs-sidecar/pkg/event"
	"github.com/tangxusc/cqrs-sidecar/pkg/proxy"
	"regexp"
	"strings"
)

func init() {
	variables := &ackEvent{}
	compile, e := regexp.Compile(`(?i).*\s*call ack\('(.+)'\)$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	proxy.Handlers = append(proxy.Handlers, variables)
}

type ackEvent struct {
	compile *regexp.Regexp
}

func (s *ackEvent) Match(sql string) bool {
	return s.compile.MatchString(sql)
}

func (s *ackEvent) Handler(query string, handler *proxy.ConnHandler) (*mysql.Result, error) {
	allString := s.compile.FindAllString(query, -1)
	eventId := allString[1]
	if handler.Tx == nil {
		return nil, fmt.Errorf(`请开启事务再进行ack`)
	}
	if len(strings.TrimSpace(eventId)) <= 0 {
		return nil, fmt.Errorf(`eventId不能为空`)
	}
	e := db.ConnInstance.ExecWithTx(handler.Tx, `update event set status = ? where id= ? and status= ? `, event.Confirmed, eventId, event.NotConfirmed)
	return nil, e
}
