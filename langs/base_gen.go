package langs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/samuel/go-thrift/parser"
)

type BaseGen struct {
	Lang      string
	Namespace string
	Thrifts   map[string]*parser.Thrift
}

func (g *BaseGen) Init(lang string, parsedThrift map[string]*parser.Thrift) {
	g.Lang = lang
	g.Thrifts = parsedThrift
	g.CheckNamespace()

	if err := g.checkMethodName(); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}

func (g *BaseGen) CheckNamespace() {
	for _, thrift := range g.Thrifts {
		for lang, namepace := range thrift.Namespaces {
			if lang == g.Lang {
				g.Namespace = namepace
				return
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Namespace not found for: %s\n", g.Lang)
	os.Exit(2)
}

func (g *BaseGen) checkMethodName() error {
	for _, thrift := range g.Thrifts {
		for _, s := range thrift.Services {
			for _, m := range s.Methods {
				if m.Name[:1] == strings.ToUpper(m.Name[:1]) {
					continue
				}

				return errors.New("the first letter of thrift service methods should be capitial")
			}
		}
	}

	return nil
}
