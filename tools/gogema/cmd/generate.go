// Package cmd provides the command-line interface for gogema.
package cmd

import (
	"github.com/fahmihidayah/gogema/internal/reader"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command that processes YAML definitions
// and outputs generated code based on the specified framework.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from YAML project and model definitions",
	Long: `Generate code by reading project.yml and model definitions from the
specified path directory.

The command expects the following structure:
  <path>/
    project.yml      # Project configuration (name, package, version, etc.)
    model/           # Directory containing model YAML files
      user.yml       # Example model definition
      post.yml       # Example model definition

Each model YAML file can define:
  - Fields with types, constraints, and validation rules
  - Database indexes (single and composite)
  - Relationships (has_one, has_many, belongs_to, many2many)
  - Foreign key constraints

Examples:
  gogema generate                          # Use current directory
  gogema generate --path ./my-project      # Specify project path
  gogema generate --framework golang       # Specify target framework`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		framework, _ := cmd.Flags().GetString("framework")
		project, models := reader.LoadProjectAndModel(path)
		if project == nil || models == nil {
			return
		}

		if framework == "" {
			framework = "golang"
		}

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateCmd.PersistentFlags().String("path", ".", "Path to project configuration directory")
	generateCmd.PersistentFlags().String("framework", "golang", "Output directory for generated models")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
