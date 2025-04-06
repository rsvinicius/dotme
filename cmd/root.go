package cmd

import (
	"fmt"
	"os"

	"github.com/rsvinicius/dotme/internal"
	"github.com/spf13/cobra"
)

var (
	// Version info set by the build process
	version string
	commit  string
	date    string
)

var rootCmd = &cobra.Command{
	Use:   "dotme [git-repository-url]",
	Short: "Apply dotfiles from a Git repository",
	Long: `dotme is a command line tool that applies dotfiles from a Git repository to your current working directory.
It only copies files and folders starting with a dot (.) from the root of the repository.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		err := internal.ProcessRepository(repoURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dotme version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// SetVersionInfo sets the version information from main
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = version
}

func init() {
	rootCmd.AddCommand(versionCmd)
	// Here you will define your flags and configuration settings
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
