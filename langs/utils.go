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

func (u *Util) IsNormalMethod(m *parser.Method) bool {
	return !u.IsBroadcastMethod(m) && !u.IsDirectMethod(m)
}

func (u *Util) IsBroadcastMethod(m *parser.Method) bool {
	return m.Oneway && strings.HasPrefix(m.Name, "On")
}

func (u *Util) HasDirectMethod(s *parser.Service) bool {
	for _, m := range s.Methods {
		if u.IsDirectMethod(m) {
			return true
		}
	}
	return false
}

func (u *Util) IsDirectMethod(m *parser.Method) bool {
	return strings.HasPrefix(m.Name, "Direct")
}
