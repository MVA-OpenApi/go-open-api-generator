package cli

import (
	"fmt"
	"os"
	"strconv"

	gen "go-open-api-generator/generator"
	par "go-open-api-generator/parser"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generator [command] [flags]",
	Short: "Create server and client API code from OpenApi Spec",
	Long: "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
}

var generateCmd = &cobra.Command{
	Use:   "generate [port]",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		port := args[0]

		// check if port is an int
		portNumber, err := strconv.Atoi(port);
		if  err != nil {
			fmt.Printf("Given port has to be an integer.")
			return
		}

		gen.CreateBuildDirectory()
		gen.GenerateServerTemplate(int16(portNumber))
	},
}

var parseCmd = &cobra.Command {
	Use: "parse [OpenAPI spec path] [template path]",
	Short: "Parse the input files for later use.",
	Long: "Parse the OpenAPI Specifiaction file (JSON Format) and the template file",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		spec_path := args[0]
		template_path := args[1]
		
		// check if spec file exists
		if !CheckIfFileExists(spec_path) {
			fmt.Println("Specification file doesn't exists.")
			return
		}

		// check if template file exists
		if !CheckIfFileExists(template_path) {
			fmt.Println("Template file doesn't exists.")
			return
		}
		
		// call parser function
		par.Parse(spec_path, template_path)
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
	rootCmd.AddCommand(generateCmd, parseCmd)
}
