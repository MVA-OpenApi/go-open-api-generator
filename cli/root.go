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
var (
	projectPath  string
	projectName  string
	loggerFlag   bool
	databaseFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-open-api-generator",
	Short: "Create server and client API code from OpenApi Spec",
	Long:  "Generate Go-Server code and RapidDoc-Clientcode for your application by providing an OpenAPI Specification",
}

var generateCmd = &cobra.Command{
	Use:     "generate <path to OpenAPI Specification>",
	Short:   "Create server and client API code from OpenApi Spec",
	Long:    "Generate Go-Server code and RapidDoc-Clientcode for your application by providing an OpenAPI Specification",
	Example: "generate ./stores.yaml -o ./outputPath -n StoresAPI",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		openAPIPath := args[0]

		// output project path
		if projectPath == "" {
			projectPath, _ = os.UserHomeDir()
		}

		// output project name
		if projectName == "" {
			projectName = "build"
		}

		projectDestination := filepath.Join(projectPath, projectName)
		config := gen.GeneratorConfig{
			OpenAPIPath:  openAPIPath,
			OutputPath:   projectDestination,
			ModuleName:   projectName,
			DatabaseName: "database",
			Flags: gen.Flags{
				UseDatabase: databaseFlag,
				UseLogger:   loggerFlag,
			},
		}

		log.Info().Msg("Generating project...")
		err := gen.GenerateServer(config)

		if err != nil {
			log.Error().Msg("Aborting...")
			return
		}

		log.Info().Msg("Running external commands...")
		log.Info().Msg("RUN `go mod init " + projectName + "`")
		extCmd.RunCommand("go mod init "+projectName, projectDestination)
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
	generateCmd.Flags().StringVarP(&projectPath, "output", "o", "", "path where generated code gets stored (default is the home directory)")
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "", "module name of generated code (default is 'build')")
	generateCmd.Flags().BoolVarP(&loggerFlag, "logger", "l", false, "use logging middleware in generated code (default is 'false')")
	generateCmd.Flags().BoolVarP(&databaseFlag, "database", "d", false, "add sqlite3 database in generated code (default is 'false')")

	// add generate command
	rootCmd.AddCommand(generateCmd)
}
