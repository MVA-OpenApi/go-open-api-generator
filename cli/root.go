package cli

import (
	"fmt"
	gen "go-open-api-generator/generator"
	"os"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generator [command] [flags]",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
}

var generateCmd = &cobra.Command{
	Use:   "generate [port]",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		input_path := args[0]

		if !CheckIfFileExists(input_path) {
			fmt.Println("No valid input file path given.")
			return
		}

		loader := openapi3.NewLoader()

		spec, err := loader.LoadFromFile(input_path)
		if err != nil {
			fmt.Println("Error loading File", err)
			return
		}

		err = spec.Validate(loader.Context)
		if err != nil {
			fmt.Println("Not a valid OpenAPI Spec!")
			fmt.Println(err)
			return
		}

		// TODO use proper logging

		fmt.Printf("Loaded Spec \"%s\" (Version %s)\n", spec.Info.Title, spec.Info.Version)
		fmt.Println("With available operations: ")
		for path, path_item := range spec.Paths {
			for op_string, op := range path_item.Operations() {
				fmt.Printf("%s %s: %s\n", op_string, path, op.Summary)
			}
		}

		port, _ := strconv.Atoi(spec.Servers[0].Variables["port"].Default)

		gen.CreateBuildDirectory()
		gen.GenerateServerTemplate(int16(port))
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
