package langs

import (
	"strings"

	"github.com/samuel/go-thrift/parser"
)

func IsBroadcastMethod(m *parser.Method) bool {
	return m.Oneway && strings.HasPrefix(m.Name, "On")
}
