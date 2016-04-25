package csharp

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Wuvist/go-thrift/parser"
	"github.com/ezbuy/ezrpc/langs"
)

var lang = "csharp"

type gen struct {
	langs.BaseGen
}

func (g *gen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	g.BaseGen.Init(lang, parsedThrift)

	path, err := filepath.Abs(output)
	if err != nil {
		log.Fatalf("failed to get absolute path of [%s]", output)
	}

	for _, parsed := range parsedThrift {
		ns := getNamespace(parsed)

		// make output dir
		d := filepath.Join(path, filepath.Join(strings.Split(ns, ".")...))
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalf("failed to make output dir [%s]", d)
		}

		// write file
		for n, s := range parsed.Services {
			f := filepath.Join(d, "gen_"+n+"_server.cs")
			data := &tmpldata{Namespace: ns, Service: s}

			if err := writefile(f, data); err != nil {
				log.Fatalf("failed to write file [%s]: %s", f, err.Error())
			}
		}
	}
}

func getNamespace(t *parser.Thrift) string {
	if namespace, ok := t.Namespaces[lang]; ok {
		return namespace
	}
	return ""
}

func init() {
	langs.Langs[lang] = &gen{}
}
