package cmd

import (
	"fmt"
	"os"

	"github.com/rsvinicius/dotme/internal"
	"github.com/rsvinicius/dotme/internal/alias"
	"github.com/rsvinicius/dotme/internal/patterns"
	"github.com/spf13/cobra"
)

var (
	// Version info set by the build process
	version string
	commit  string
	date    string

	// Command flags
	aliasFlag       string
	saveFlag        string
	includePatterns string
	excludePatterns string
)

var rootCmd = &cobra.Command{
	Use:   "dotme [git-repository-url]",
	Short: "Apply dotfiles from a Git repository",
	Long: `dotme is a command line tool that applies dotfiles from a Git repository to your current working directory.
It only copies files and folders starting with a dot (.) from the root of the repository.

You can use include and exclude patterns to filter which dotfiles are copied:
  --include: Comma-separated list of patterns to include (e.g., ".vscode,.gitconfig")
  --exclude: Comma-separated list of patterns to exclude (e.g., ".DS_Store")

Patterns support glob matching (*, ?, [abc], etc.).`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse patterns
		includeList := patterns.ParsePatterns(includePatterns)
		excludeList := patterns.ParsePatterns(excludePatterns)

		// Check for alias flag first
		if aliasFlag != "" {
			repoURL, err := alias.GetRepo(aliasFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("🔍 Using alias '%s' for repository: %s\n", aliasFlag, repoURL)
			err = internal.ProcessRepository(repoURL, includeList, excludeList)
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
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying help: %s\n", err)
			}
			os.Exit(1)
		}

		repoURL := args[0]
		err := internal.ProcessRepository(repoURL, includeList, excludeList)
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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
}

var setDefaultPatternsCmd = &cobra.Command{
	Use:   "set-default-patterns",
	Short: "Set default include/exclude patterns",
	Long: `Set default include and exclude patterns that will be used when no patterns are specified.

Examples:
  dotme config set-default-patterns --include=".vscode,.gitconfig" --exclude=".DS_Store"
  dotme config set-default-patterns --include=".git*"
  dotme config set-default-patterns --exclude=".DS_Store,.Trash*"`,
	Run: func(cmd *cobra.Command, args []string) {
		includeList := patterns.ParsePatterns(includePatterns)
		excludeList := patterns.ParsePatterns(excludePatterns)

		patternConfig := alias.PatternConfig{
			IncludePatterns: includeList,
			ExcludePatterns: excludeList,
		}

		err := alias.SetDefaultPatterns(patternConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Default patterns updated successfully")
		if len(includeList) > 0 {
			fmt.Printf("   Include patterns: %v\n", includeList)
		}
		if len(excludeList) > 0 {
			fmt.Printf("   Exclude patterns: %v\n", excludeList)
		}
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// Show aliases
		aliases, err := alias.ListAliases()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading aliases: %s\n", err)
		} else {
			fmt.Println("📋 Repository aliases:")
			if len(aliases) == 0 {
				fmt.Println("   (none)")
			} else {
				for name, url := range aliases {
					fmt.Printf("   📎 %s: %s\n", name, url)
				}
			}
		}

		// Show default patterns
		defaultPatterns, err := alias.GetDefaultPatterns()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading default patterns: %s\n", err)
		} else {
			fmt.Println("\n🔍 Default patterns:")
			if len(defaultPatterns.IncludePatterns) == 0 && len(defaultPatterns.ExcludePatterns) == 0 {
				fmt.Println("   (none - will include all dotfiles)")
			} else {
				if len(defaultPatterns.IncludePatterns) > 0 {
					fmt.Printf("   Include: %v\n", defaultPatterns.IncludePatterns)
				}
				if len(defaultPatterns.ExcludePatterns) > 0 {
					fmt.Printf("   Exclude: %v\n", defaultPatterns.ExcludePatterns)
				}
			}
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
	rootCmd.AddCommand(configCmd)

	// Add config subcommands
	configCmd.AddCommand(setDefaultPatternsCmd)
	configCmd.AddCommand(showConfigCmd)

	// Add flags to root command
	rootCmd.Flags().StringVarP(&aliasFlag, "alias", "a", "", "Use a saved repository by alias")
	rootCmd.Flags().StringVarP(&saveFlag, "save", "s", "", "Save the repository with the given alias")
	rootCmd.Flags().StringVar(&includePatterns, "include", "", "Comma-separated list of patterns to include (e.g., '.vscode,.gitconfig')")
	rootCmd.Flags().StringVar(&excludePatterns, "exclude", "", "Comma-separated list of patterns to exclude (e.g., '.DS_Store')")

	// Add flags to set-default-patterns command
	setDefaultPatternsCmd.Flags().StringVar(&includePatterns, "include", "", "Comma-separated list of default include patterns")
	setDefaultPatternsCmd.Flags().StringVar(&excludePatterns, "exclude", "", "Comma-separated list of default exclude patterns")
}