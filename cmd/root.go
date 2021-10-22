package cmd

import (
	"fmt"

	"github.com/mmiranda/markdown-index/markdown"
	"github.com/spf13/cobra"
)

var (
	directory, output, skipDirectory, useHeading string
	// maxDepth          int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "markdown-index",
	Short: "Generate summary index for your Markdown files",
	Long: `markdown-index iterates through a given directory, 
reads all markdown files recursively and generate for you 
a index file with the summary of each file found.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing script...")
		markdown.Execute(output, directory)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version numbe`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1-0-alpha")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVar(&directory, "directory", "./", "Directory to search markdown files recursively")
	rootCmd.PersistentFlags().StringVar(&output, "output", "markdown-index.md", "Final markdown file to be created")
	rootCmd.PersistentFlags().StringVar(&skipDirectory, "skipDirectory", "", "Skip directory in the recursive walk")
	rootCmd.PersistentFlags().StringVar(&useHeading, "heading", "", "Use this Heading inside the markdown file as summary of the file")
	// rootCmd.PersistentFlags().StringVar(&output, "maxDepth", "5", "Maximum depth level to look for files")

}
