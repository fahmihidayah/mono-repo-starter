// Package cmd provides the command-line interface for gogema.
// Gogema is a code generation tool that reads YAML-based project and model
// definitions to generate boilerplate code for Go applications.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogema",
	Short: "A YAML-driven code generator for Go applications",
	Long: `Gogema is a code generation CLI tool that transforms YAML-based
project and model definitions into Go source code.

Define your project configuration in project.yml and your data models
in the model/ directory as YAML files. Gogema reads these definitions
and generates the corresponding Go structs, database models, and more.

Features:
  - YAML-based project and model configuration
  - Support for fields, indexes, relationships, and foreign keys
  - Customizable field attributes (validation, defaults, constraints)
  - Extensible framework support

Example usage:
  gogema generate --path ./my-project
  gogema generate --path ./my-project --framework golang`,
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogema.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
