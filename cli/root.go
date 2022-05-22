package cli

import (
	extCmd "go-open-api-generator/cmd"
	gen "go-open-api-generator/generator"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// variables for the flags
var projectPath string
var projectName string
var openAPIPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generator [flags]",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
}

var generateCmd = &cobra.Command{
	Use:   "generate -o <path to OpenAPI Specification>",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and ReactJS-Clientcode for your application by providing an OpenAPI Specification",
	Run: func(cmd *cobra.Command, args []string) {
		// output project path
		if projectPath == "" {
			projectPath,_ = os.UserHomeDir()
		}

		// output project name
		if projectName == "" {
			projectName = "build"
		}

		// template path
		if openAPIPath == "" {
			log.Error().Msg("No OpenAPI path given, add -o <OpenAPI Specificatio path> flag.")
			return
		}		
		
		projectDestination := filepath.Join(projectPath, projectName) 

		log.Info().Msg("Generating project...")
		gen.GenerateServer(openAPIPath, projectDestination, projectName)

		log.Info().Msg("Running external commands...")
		log.Info().Msg("RUN `go mod init " + projectName + "`")
		extCmd.RunCommand("go mod init " + projectName, projectDestination)
		log.Info().Msg("RUN `go mod tidy`")
		extCmd.RunCommand("go mod tidy", projectDestination)
		log.Info().Msg("RUN `go fmt ./...`")
		extCmd.RunCommand("go fmt ./...", projectDestination)
		log.Info().Msg("DONE project created at: " + projectDestination)
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
	// add generate flags
	generateCmd.Flags().StringVarP(&projectPath, "output project path", "p", "", "path where generated code gets stored")
	generateCmd.Flags().StringVarP(&projectName, "name of the generated project", "n", "", "module name of generated code")
	generateCmd.Flags().StringVarP(&openAPIPath, "OpenAPI path", "o", "", "path where the OpenAPI Specification is stored")

	// add generate command
	rootCmd.AddCommand(generateCmd)
}

