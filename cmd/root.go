package cmd

import (
	"fmt"
	"os"

	"github.com/rsvinicius/dotme/internal"
	"github.com/rsvinicius/dotme/internal/alias"
	"github.com/spf13/cobra"
)

var (
	// Version info set by the build process
	version string
	commit  string
	date    string

	// Command flags
	aliasFlag string
	saveFlag  string
)

var rootCmd = &cobra.Command{
	Use:   "dotme [git-repository-url]",
	Short: "Apply dotfiles from a Git repository",
	Long: `dotme is a command line tool that applies dotfiles from a Git repository to your current working directory.
It only copies files and folders starting with a dot (.) from the root of the repository.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check for alias flag first
		if aliasFlag != "" {
			repoURL, err := alias.GetRepo(aliasFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("🔍 Using alias '%s' for repository: %s\n", aliasFlag, repoURL)
			err = internal.ProcessRepository(repoURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			return
		}

		// Check for save flag
		if saveFlag != "" {
			if len(args) != 1 {
				fmt.Fprintf(os.Stderr, "Error: repository URL is required when using --save\n")
				os.Exit(1)
			}
			repoURL := args[0]
			err := alias.SaveRepo(repoURL, saveFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Repository '%s' saved with alias '%s'\n", repoURL, saveFlag)
			return
		}

		// Normal operation - apply dotfiles from repository
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Error: repository URL is required\n")
			cmd.Help()
			os.Exit(1)
		}

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

var listCmd = &cobra.Command{
	Use:   "list-aliases",
	Short: "List all saved repository aliases",
	Run: func(cmd *cobra.Command, args []string) {
		aliases, err := alias.ListAliases()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found. Save one with 'dotme -s <alias> <repository-url>'")
			return
		}

		fmt.Println("📋 Saved repository aliases:")
		fmt.Println("----------------------------")
		for name, url := range aliases {
			fmt.Printf("📎 %s: %s\n", name, url)
		}
	},
	Aliases: []string{"ls"},
}

var removeAliasCmd = &cobra.Command{
	Use:   "remove-alias <alias-name>",
	Short: "Remove a saved repository alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		aliasName := args[0]
		err := alias.DeleteAlias(aliasName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Alias '%s' removed successfully\n", aliasName)
	},
	Aliases: []string{"rm"},
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
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeAliasCmd)

	// Add flags
	rootCmd.Flags().StringVarP(&aliasFlag, "alias", "a", "", "Use a saved repository by alias")
	rootCmd.Flags().StringVarP(&saveFlag, "save", "s", "", "Save the repository with the given alias")
}
