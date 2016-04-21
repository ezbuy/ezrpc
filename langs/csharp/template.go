package csharp

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/ezbuy/ezrpc/langs"
	"github.com/ezbuy/ezrpc/tmpl"
	"github.com/samuel/go-thrift/parser"
)

var tpl *template.Template

type tmpldata struct {
	Namespace string
	Service   *parser.Service
}

func writefile(f string, d *tmpldata) error {
	file, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return tpl.ExecuteTemplate(file, "ezrpc/csharp", d)
}

func init() {
	tpl = template.New("ezrpc/csharp")

	funcs := template.FuncMap{
		"ToLower": strings.ToLower,
		"Utils":   langs.Utils,
	}
	tpl.Funcs(funcs)

	files := []string{
		"tmpl/csharp/server.gocs",
	}

	for _, f := range files {
		data, err := tmpl.Asset(f)
		if err != nil {
			log.Fatalf("failed to get template [%s]: %s", f, err.Error())
		}

		if _, err := tpl.Parse(string(data)); err != nil {
			log.Fatalf("failed to parse template [%s]: %s", f, err.Error())
		}
	}
}
