package cmd

import (
	"fmt"
	"os"

	"github.com/rsvinicius/dotme/internal"
	"github.com/spf13/cobra"
)

const (
	// Version is the current version of dotme
	Version = "v0.1.0"
)

var rootCmd = &cobra.Command{
	Use:     "dotme [git-repository-url]",
	Short:   "Apply dotfiles from a Git repository",
	Long: `dotme is a command line tool that applies dotfiles from a Git repository to your current working directory.
It only copies files and folders starting with a dot (.) from the root of the repository.`,
	Version: Version,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		err := internal.ProcessRepository(repoURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
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

func init() {
	// Here you will define your flags and configuration settings
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
