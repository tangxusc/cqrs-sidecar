package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-sidecar/pkg/db"
	"github.com/tangxusc/cqrs-sidecar/pkg/proxy"
	"regexp"
)

func init() {
	set := &set{}
	compile, e := regexp.Compile(`(?i).*\s*set\s*.*$`)
	if e != nil {
		panic(e.Error())
	}
	set.compile = compile
	proxy.Handlers = append(proxy.Handlers, set)
}

type set struct {
	compile *regexp.Regexp
}

func (s *set) Match(sql string) bool {
	return s.compile.MatchString(sql)
}

func (s *set) Handler(query string, handler *proxy.ConnHandler) (*mysql.Result, error) {
	_, _, err := db.ConnInstance.Proxy(query)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
