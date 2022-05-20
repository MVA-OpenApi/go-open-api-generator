package cli

import (
	gen "go-open-api-generator/generator"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generator [command] [flags]",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
}

var generateCmd = &cobra.Command{
	Use:   "generate [open-api-file-path]",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gen.GenerateServer(args[0])
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// add sub commands
	rootCmd.AddCommand(generateCmd)
}
