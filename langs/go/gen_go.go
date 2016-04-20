package gogen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ezbuy/ezrpc/langs"
	"github.com/samuel/go-thrift/parser"
)

const langName = "go"

type GoGen struct {
	langs.BaseGen
}

func getNamespace(namespaces map[string]string) string {
	if namespace, ok := namespaces[langName]; ok {
		return namespace
	}

	return ""
}

func genNamespace(namespace string) (string, string) {
	path := strings.Replace(namespace, ".", "/", -1)
	pkgName := filepath.Base(path)
	return path, pkgName
}

func panicWithErr(format string, msg ...interface{}) {
	panic(fmt.Errorf(format, msg...))
}

type ServerData struct {
	Namespace string
	Service   *parser.Service
}

func (d ServerData) HasBroadcastMethod() bool {
	for _, m := range d.Service.Methods {
		if langs.IsBroadcastMethod(m) {
			return true
		}
	}

	return false
}

func (this *GoGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	this.BaseGen.Init(langName, parsedThrift)

	outputPath, err := filepath.Abs(output)
	if err != nil {
		panicWithErr("fail to get absolute path for %q", output)
	}

	outputPackageDirs := make([]string, 0, len(parsedThrift))

	fmt.Println("##### Parsing:")
	for filename, parsed := range parsedThrift {
		fmt.Printf("%s\n", filename)
		namespace := getNamespace(parsed.Namespaces)
		importPath, _ := genNamespace(namespace)

		// make output dir
		pkgDir := filepath.Join(outputPath, importPath)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			panicWithErr("fail to make package directory %s", pkgDir)
		}

		outputPackageDirs = append(outputPackageDirs, pkgDir)

		// write file
		for name, service := range parsed.Services {
			fname := filepath.Join(pkgDir, "gen_"+name+"_server.go")
			data := ServerData{
				Namespace: namespace,
				Service:   service,
			}
			if err := outputFile(fname, "server", data); err != nil {
				panicWithErr("fail to write defines file %q : %s", fname, err)
			}
		}

	}
}

func init() {
	langs.Langs[langName] = &GoGen{}
}
