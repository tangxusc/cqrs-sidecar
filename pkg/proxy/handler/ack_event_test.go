package handler

import (
	"fmt"
	"regexp"
	"testing"
)

func TestAHandler(t *testing.T) {
	compile, e := regexp.Compile(`(?i).*\s*call ack\('(.+)'\)$`)
	if e != nil {
		panic(e.Error())
	}
	s := `call ack('1')`
	allString := compile.FindAllString(s, -1)
	fmt.Println(allString)
	submatch := compile.FindAllStringSubmatch(s, -1)
	println(submatch[0][1])
	fmt.Println(submatch)
}
