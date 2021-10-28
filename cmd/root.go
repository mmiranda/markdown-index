package cmd

import (
	"fmt"

	"github.com/mmiranda/markdown-index/markdown"
	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
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
	Run: func(cmd *cobra.Command, args []string) {
		markdown.Execute(output, directory)
	},
}

var (
	shortened     = false
	version       = "dev"
	commit        = "none"
	date          = "unknown"
	versionOutput = "json"
	versionCmd    = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			resp := goVersion.FuncWithOutput(shortened, version, commit, date, versionOutput)
			fmt.Print(resp)
			return
		},
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	versionCmd.Flags().StringVarP(&output, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")
	rootCmd.AddCommand(versionCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVar(&directory, "directory", ".", "Directory to search markdown files recursively")
	rootCmd.PersistentFlags().StringVar(&output, "output", "markdown-index.md", "Final markdown file to be created")
	rootCmd.PersistentFlags().StringVar(&skipDirectory, "skipDirectory", "", "Skip directory in the recursive walk")
	rootCmd.PersistentFlags().StringVar(&useHeading, "heading", "", "Use this Heading inside the markdown file as summary of the file")
	// rootCmd.PersistentFlags().StringVar(&output, "maxDepth", "5", "Maximum depth level to look for files")

}
