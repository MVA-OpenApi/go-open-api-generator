package cli

import (
	"fmt"
	"os"

	gen "go-open-api-generator/generator"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "genoapi generate [input file path] [flags]",
	Short: "Create server and client API code from OpenApi Spec",
	Long: "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
}

var generateCmd = &cobra.Command{
	Use:   "generate [input file path]",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		
		input_path := args[0]
		
		if !CheckIfFileExists(input_path) {
			fmt.Println("No valid input file path given.")
			return
		}
		
		jsonFile, err := os.Open(input_path)
		
		if err != nil {
			fmt.Println(err)
		}
		
		defer jsonFile.Close()

		// TODO parse OPenAPI spec

		gen.CreateBuildDirectory()
		gen.GenerateServerTemplate(3000)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.example.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// add sub commands
	rootCmd.AddCommand(generateCmd)
}
