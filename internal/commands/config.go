package commands

import (
	"fmt"
	"strings"

	"github.com/agarcher/wt/internal/config"
	"github.com/agarcher/wt/internal/userconfig"
	"github.com/spf13/cobra"
)

var (
	configGlobal     bool
	configUnset      bool
	configList       bool
	configShowOrigin bool
)

func init() {
	configCmd.Flags().BoolVar(&configGlobal, "global", false, "Set/get global configuration")
	configCmd.Flags().BoolVar(&configUnset, "unset", false, "Remove a per-repo configuration value")
	configCmd.Flags().BoolVar(&configList, "list", false, "List all configuration values")
	configCmd.Flags().BoolVar(&configShowOrigin, "show-origin", false, "Show where each configuration value comes from")
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Manage user configuration",
	Long: `Get and set wt user configuration options.

User settings are stored in ~/.config/wt/config.yaml

Configuration keys:
  remote    Remote to compare against (empty = local comparison)
  fetch     Auto-fetch before list/cleanup (only applies when remote is set)

Examples:
  wt config --list                  # List all settings
  wt config --show-origin           # Show where each value comes from
  wt config fetch                   # Get the value of 'fetch'
  wt config --global remote origin  # Set global remote
  wt config --global fetch true     # Enable auto-fetch globally
  wt config remote upstream         # Set remote for current repo only
  wt config --unset remote          # Remove per-repo remote override

Note: 'fetch' only has an effect when 'remote' is set. If remote is empty,
comparisons are done against the local branch and fetch is ignored.`,
	RunE: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	// Load user config
	cfg, err := userconfig.Load()
	if err != nil {
		return fmt.Errorf("failed to load user config: %w", err)
	}

	// Handle --list
	if configList {
		return printConfigList(cmd, cfg)
	}

	// Handle --show-origin
	if configShowOrigin {
		return printConfigShowOrigin(cmd, cfg)
	}

	// Handle --unset
	if configUnset {
		if len(args) < 1 {
			return fmt.Errorf("usage: wt config --unset <key>")
		}
		return unsetConfig(cmd, cfg, args[0])
	}

	// Get or set
	switch len(args) {
	case 0:
		return fmt.Errorf("usage: wt config [--global] <key> [value]\n       wt config --list\n       wt config --show-origin")
	case 1:
		// Get value
		return getConfig(cmd, cfg, args[0])
	case 2:
		// Set value
		return setConfig(cmd, cfg, args[0], args[1])
	default:
		return fmt.Errorf("too many arguments")
	}
}

func printConfigList(cmd *cobra.Command, cfg *userconfig.UserConfig) error {
	out := cmd.OutOrStdout()

	// Print global values
	if cfg.Remote != "" {
		_, _ = fmt.Fprintf(out, "remote = %s (global)\n", cfg.Remote)
	}
	if cfg.Fetch {
		_, _ = fmt.Fprintf(out, "fetch = true (global)\n")
	} else if cfg.Remote != "" {
		// Only show fetch=false if remote is set (otherwise it's meaningless)
		_, _ = fmt.Fprintf(out, "fetch = false (global)\n")
	}

	// Print per-repo values
	for repoPath, repoConfig := range cfg.Repos {
		if repoConfig.Remote != "" {
			_, _ = fmt.Fprintf(out, "repos.%s.remote = %s\n", repoPath, repoConfig.Remote)
		}
		if repoConfig.Fetch != nil {
			_, _ = fmt.Fprintf(out, "repos.%s.fetch = %v\n", repoPath, *repoConfig.Fetch)
		}
	}

	return nil
}

func printConfigShowOrigin(cmd *cobra.Command, cfg *userconfig.UserConfig) error {
	out := cmd.OutOrStdout()

	// Get current repo path for context
	repoRoot, _ := config.GetMainRepoRoot()

	configPath, err := userconfig.GetConfigPath()
	if err != nil {
		configPath = "(unknown)"
	}

	// Show effective values for current repo
	if repoRoot != "" {
		remote := cfg.GetRemoteForRepo(repoRoot)
		fetch := cfg.GetFetchForRepo(repoRoot)

		// Determine source of remote
		if repoConfig, ok := cfg.Repos[repoRoot]; ok && repoConfig.Remote != "" {
			_, _ = fmt.Fprintf(out, "remote = %-20s %s (repos.%s)\n", remote, configPath, repoRoot)
		} else if cfg.Remote != "" {
			_, _ = fmt.Fprintf(out, "remote = %-20s %s (global)\n", remote, configPath)
		} else {
			_, _ = fmt.Fprintf(out, "remote = %-20s (default: local comparison)\n", "\"\"")
		}

		// Determine source of fetch
		if repoConfig, ok := cfg.Repos[repoRoot]; ok && repoConfig.Fetch != nil {
			_, _ = fmt.Fprintf(out, "fetch = %-21v %s (repos.%s)\n", fetch, configPath, repoRoot)
		} else if cfg.Fetch {
			_, _ = fmt.Fprintf(out, "fetch = %-21v %s (global)\n", fetch, configPath)
		} else {
			_, _ = fmt.Fprintf(out, "fetch = %-21v (default)\n", false)
		}

		// Show repo's default_branch if set
		if repoCfg, err := config.Load(repoRoot); err == nil && repoCfg.DefaultBranch != "" {
			_, _ = fmt.Fprintf(out, "default_branch = %-14s .wt.yaml (repo)\n", repoCfg.DefaultBranch)
		}
	} else {
		// Not in a repo, just show global values
		_, _ = fmt.Fprintf(out, "remote = %-20s %s (global)\n", cfg.Remote, configPath)
		_, _ = fmt.Fprintf(out, "fetch = %-21v %s (global)\n", cfg.Fetch, configPath)
	}

	return nil
}

func getConfig(cmd *cobra.Command, cfg *userconfig.UserConfig, key string) error {
	// Validate key
	if !isValidKey(key) {
		return fmt.Errorf("unknown config key: %s\nValid keys: %s", key, strings.Join(userconfig.ValidKeys(), ", "))
	}

	if configGlobal {
		// Get global value
		value, err := cfg.GetGlobal(key)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), value)
	} else {
		// Get effective value for current repo
		repoRoot, err := config.GetMainRepoRoot()
		if err != nil {
			return fmt.Errorf("not in a git repository (use --global for global config)")
		}

		switch key {
		case "remote":
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), cfg.GetRemoteForRepo(repoRoot))
		case "fetch":
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), cfg.GetFetchForRepo(repoRoot))
		}
	}

	return nil
}

func setConfig(cmd *cobra.Command, cfg *userconfig.UserConfig, key, value string) error {
	// Validate key
	if !isValidKey(key) {
		return fmt.Errorf("unknown config key: %s\nValid keys: %s", key, strings.Join(userconfig.ValidKeys(), ", "))
	}

	// Validate fetch value
	if key == "fetch" && value != "true" && value != "false" {
		return fmt.Errorf("fetch must be 'true' or 'false'")
	}

	if configGlobal {
		// Set global value
		if err := cfg.SetGlobal(key, value); err != nil {
			return err
		}

		// Warn if setting fetch=true without remote
		if key == "fetch" && value == "true" && cfg.Remote == "" {
			cmd.PrintErrln("Warning: fetch has no effect when remote is not set")
		}
	} else {
		// Set per-repo value
		repoRoot, err := config.GetMainRepoRoot()
		if err != nil {
			return fmt.Errorf("not in a git repository (use --global for global config)")
		}

		if err := cfg.SetForRepo(repoRoot, key, value); err != nil {
			return err
		}

		// Warn if setting fetch=true without remote for this repo
		if key == "fetch" && value == "true" && cfg.GetRemoteForRepo(repoRoot) == "" {
			cmd.PrintErrln("Warning: fetch has no effect when remote is not set")
		}
	}

	// Save config
	if err := userconfig.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func unsetConfig(cmd *cobra.Command, cfg *userconfig.UserConfig, key string) error {
	// Validate key
	if !isValidKey(key) {
		return fmt.Errorf("unknown config key: %s\nValid keys: %s", key, strings.Join(userconfig.ValidKeys(), ", "))
	}

	if configGlobal {
		// Unset global value
		if err := cfg.UnsetGlobal(key); err != nil {
			return err
		}
	} else {
		// Get current repo
		repoRoot, err := config.GetMainRepoRoot()
		if err != nil {
			return fmt.Errorf("not in a git repository (use --global to unset global config)")
		}

		// Unset per-repo value
		if err := cfg.UnsetForRepo(repoRoot, key); err != nil {
			return err
		}
	}

	// Save config
	if err := userconfig.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func isValidKey(key string) bool {
	for _, k := range userconfig.ValidKeys() {
		if k == key {
			return true
		}
	}
	return false
}
