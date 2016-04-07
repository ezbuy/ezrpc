package csharp

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ezbuy/ezrpc/langs"
	"github.com/samuel/go-thrift/parser"
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

	log.Println(path)

	for file, parsed := range parsedThrift {
		log.Println(file)
		log.Printf("%+v", parsed)

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

	// outputPackageDirs := make([]string, 0, len(parsedThrift))

	// fmt.Println("##### Parsing:")
	// for filename, parsed := range parsedThrift {
	// 	fmt.Printf("%s\n", filename)
	// 	namespace := getNamespace(parsed.Namespaces)
	// 	importPath, _ := genNamespace(namespace)

	// 	// make output dir
	// 	pkgDir := filepath.Join(outputPath, importPath)
	// 	if err := os.MkdirAll(pkgDir, 0755); err != nil {
	// 		panicWithErr("fail to make package directory %s", pkgDir)
	// 	}

	// 	outputPackageDirs = append(outputPackageDirs, pkgDir)

	// 	// write file
	// 	for name, service := range parsed.Services {
	// 		fname := filepath.Join(pkgDir, "gen_"+name+"_server.go")
	// 		data := ServerData{
	// 			Namespace: namespace,
	// 			Service:   service,
	// 		}
	// 		if err := outputFile(fname, "server", data); err != nil {
	// 			panicWithErr("fail to write defines file %q : %s", fname, err)
	// 		}
	// 	}

	// }
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
