package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Wuvist/go-thrift/parser"
	"github.com/ezbuy/ezrpc/global"
	"github.com/ezbuy/ezrpc/langs"
	_ "github.com/ezbuy/ezrpc/langs/csharp" // fullfill langs
	_ "github.com/ezbuy/ezrpc/langs/go"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate ezrpc server/client code",
	Run: func(cmd *cobra.Command, args []string) {
		if lang == "" {
			fmt.Println("-l language must be specified")
			return
		}

		if input == "" {
			fmt.Println("-i input thrift file must be specified")
			return
		}

		if output == "" {
			fmt.Println("-o output path must be specified")
			return
		}

		// initialize global variables here

		a, err := filepath.Abs(input)
		if err != nil {
			log.Fatalf("failed to get absoulte path of input file %q: %s", input, err.Error())
		}

		global.InputFile = a
		global.IsGenSrvRecursive = genSrvRecursive

		p := &parser.Parser{}
		parsedThrift, _, err := p.ParseFile(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}

		if generator, ok := langs.Langs[lang]; ok {
			generator.Generate(output, parsedThrift)
		} else {
			fmt.Printf("lang %s is not supported\n", lang)
			fmt.Println("Supported language options are:")
			for key := range langs.Langs {
				fmt.Printf("\t%s\n", key)
			}
		}
	},
}

var lang, input, output string

// when a thrift includes other thrifts, if we generate all service as server,
// we may encounter problem such as compile error due to lack of structs.
// so in the minimum, we should only generate server & structs specified to the thrift itself
var genSrvRecursive bool

func init() {
	genCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "", "language: go | csharp")
	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genCmd.PersistentFlags().BoolVarP(&genSrvRecursive, "srvRecursive", "R", true, "recursivly generate or not")

	RootCmd.AddCommand(genCmd)
}
