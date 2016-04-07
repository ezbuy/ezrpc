package cmd

import (
	"fmt"
	"os"

	"github.com/ezbuy/ezrpc/langs"
	_ "github.com/ezbuy/ezrpc/langs/csharp" // fullfill langs
	_ "github.com/ezbuy/ezrpc/langs/go"
	"github.com/samuel/go-thrift/parser"
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

func init() {
	genCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "", "language: go | csharp")
	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	RootCmd.AddCommand(genCmd)
}
