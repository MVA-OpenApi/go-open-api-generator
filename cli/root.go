package cli

import (
	extCmd "go-open-api-generator/cmd"
	gen "go-open-api-generator/generator"
	"os"

	"github.com/rs/zerolog/log"
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
		openAPIFilepath := args[0]
		projectPath := "build" // TODO: get path from CLI

		log.Info().Msg("Generating project...")
		gen.GenerateServer(openAPIFilepath)

		log.Info().Msg("Running external commands...")
		log.Info().Msg("RUN `go mod tidy`")
		extCmd.RunCommand("go mod tidy", projectPath)
		log.Info().Msg("RUN `go fmt ./...`")
		extCmd.RunCommand("go fmt ./...", projectPath)
		log.Info().Msg("DONE project created at: " + projectPath)
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
