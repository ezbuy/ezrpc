package gogen

import (
	"os"
	"strings"
	"text/template"

	"github.com/ezbuy/ezrpc/langs"
	"github.com/ezbuy/ezrpc/tmpl"
)

var tpl *template.Template

func Tpl() *template.Template {
	return tpl
}

func init() {
	tpl = template.New("ezrpc/golang")
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		"Utils":   langs.Utils,
	}
	tpl.Funcs(funcMap)
	files := []string{
		"tmpl/golang/server.gogo",
	}

	for _, filename := range files {
		data, err := tmpl.Asset(filename)
		if err != nil {
			panic(err)
		}

		if _, err = tpl.Parse(string(data)); err != nil {
			panic(err)
		}
	}
}

func outputFile(path string, tplName string, data interface{}) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return tpl.ExecuteTemplate(file, tplName, data)
}
