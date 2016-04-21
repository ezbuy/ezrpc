package langs

import (
	"strings"

	"github.com/samuel/go-thrift/parser"
)

type Util struct {
}

func Utils() *Util {
	return _util
}

var _util *Util

func (u *Util) IsBroadcastMethod(m *parser.Method) bool {
	return m.Oneway && strings.HasPrefix(m.Name, "On")
}

func (u *Util) IsDirectMethod(m *parser.Method) bool {
	return m.Oneway && strings.HasPrefix(m.Name, "Direct")
}
