package cmd

import (
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/strategy/golang"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/strategy/wasm"
	"os"

	"github.com/spf13/cobra"
	specoptions "github.com/tfadeyi/sloth-simple-comments/cmd/options/spec"
	"github.com/tfadeyi/sloth-simple-comments/internal/generate"
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/lang"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/options"
)

func specGenerateCmd() *cobra.Command {
	opts := specoptions.New()
	var outputDir string
	cmd := &cobra.Command{
		Use:           "generate",
		Short:         "Generate the Sloth definition specification from a given source code.",
		Long:          ``,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// if an argument is passed to the command
			// arg 1: should be the output directory where to store the output
			output, err := os.Getwd()
			if err != nil {
				return err
			}
			if len(args) == 1 {
				output = args[0]
			}
			outputDir = output
			return opts.Complete()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			var languageParser options.Option
			switch opts.SrcLanguage {
			case lang.Wasm:
				languageParser = wasm.Parser()
			default:
				languageParser = golang.Parser()
			}

			parser, err := parser.New(
				languageParser,
				options.Logger(&logger),
				options.Include(opts.IncludedDirs...))
			if err != nil {
				return err
			}
			service, err := parser.Parse(cmd.Context())
			if err != nil {
				return err
			}

			return generate.WriteSpecification(service, opts.StdOutput, outputDir, opts.Formats...)
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand(specGenerateCmd())
}