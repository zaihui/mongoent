package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zaihui/mongoent"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mongo-ent",
	Short: "a code generation tool for mongo model CRUD ",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ent-factory.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.PersistentFlags().StringP("schemaFile", "s", "", "file which model schema defined")
	RootCmd.PersistentFlags().StringP("outputPath", "o", "", "path to write factories")
	RootCmd.PersistentFlags().StringP("projectPath", "p", "", "the relative path of this project")
	RootCmd.PersistentFlags().StringP("goModPath", "g", "", "the relative path of this project")
}

func Fatal(msg string) {
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	Fatal(fmt.Sprintf(format, v...))
}

func ExtraFlags(cmd *cobra.Command) (string, string, string, error) {
	modelFile, err := cmd.Flags().GetString("schemaFile")
	if err != nil {
		Fatalf("get schema file failed: %v\n", err)
	}
	outputPath, err := cmd.Flags().GetString("outputPath")
	if err != nil {
		Fatalf("get output path failed: %v\n", err)
	}
	projectPath, err := cmd.Flags().GetString("projectPath")
	if err != nil {
		Fatalf("get project path failed: %v\n", err)
	}
	if projectPath == "" {
		Fatalf("project path cannot be empty")
	}
	goModPath, err := cmd.Flags().GetString("goModPath")
	if err != nil {
		Fatalf("get project path failed: %v\n", err)
	}
	if goModPath == "" {
		Fatalf("go mod path cannot be empty")
	}
	genPath := extractGenPath(outputPath, projectPath)

	return modelFile, outputPath, goModPath + genPath + "/" + mongoent.MongoSchema, err
}

func extractGenPath(outputPath, projectPath string) string {
	relativePath, err := filepath.Rel(projectPath, outputPath)
	if err != nil {
		return ""
	}
	genPath := strings.TrimPrefix(relativePath, string(filepath.Separator))
	return filepath.Join("/", genPath)
}
