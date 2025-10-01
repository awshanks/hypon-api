package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long: `Commands for viewing and managing configuration settings.
This includes viewing current configuration values and their sources.`,
}

// configViewCmd represents the config view command
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Long: `Display the current configuration values including:
- Hypon API settings
- User settings  
- Configuration file location
- Environment variable overrides`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Hypon API Configuration ===")
		fmt.Println()

		// Display configuration file information
		configFile := viper.ConfigFileUsed()
		if configFile != "" {
			fmt.Printf("Configuration file: %s\n", configFile)
		} else {
			fmt.Println("Configuration file: No config file loaded")
		}
		fmt.Println()

		// Display Hypon API settings
		fmt.Println("Hypon API Settings:")
		fmt.Printf("  api_url:  %s\n", getConfigValueWithSource("hypon.api_url", "https://api.hypon.cloud"))
		fmt.Printf("  username: %s\n", getConfigValueWithSource("hypon.username", ""))
		fmt.Printf("  password: %s\n", maskPassword(getConfigValueWithSource("hypon.password", "")))
		fmt.Printf("  oem:      %s\n", getConfigValueWithSource("hypon.oem", ""))
		fmt.Printf("  token:    %s\n", maskPassword(getConfigValueWithSource("hypon.token", "")))
		fmt.Println()

		// Display User settings
		fmt.Println("User Settings:")
		fmt.Printf("  name: %s\n", getConfigValueWithSource("user.name", ""))
		fmt.Println()

		// Display environment variable information
		fmt.Println("Environment Variables:")
		fmt.Printf("  HYPON_USER: %s\n", maskIfSet(os.Getenv("HYPON_USER")))
		fmt.Printf("  HYPON_PASS: %s\n", maskPassword(os.Getenv("HYPON_PASS")))
		fmt.Println()

		// Display config file paths that are searched
		fmt.Println("Config file search paths:")
		home, _ := os.UserHomeDir()
		fmt.Printf("  %s/.hypon-api.yaml\n", home)
		fmt.Printf("  ./.hypon-api.yaml\n")
	},
}

// configClearCmd represents the config clear command
var configClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear configuration settings",
	Long: `Clear the current configuration by removing the configuration file.
This will delete the config file and reset all settings to defaults.
Environment variables will not be affected.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		// Get the config file that would be used
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			// No config file is currently loaded, check default locations
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %v", err)
			}

			// Check common locations
			possibleConfigs := []string{
				"./.hypon-api.yaml",
				fmt.Sprintf("%s/.hypon-api.yaml", home),
			}

			// Look for existing config files
			var foundConfigs []string
			for _, path := range possibleConfigs {
				if _, err := os.Stat(path); err == nil {
					foundConfigs = append(foundConfigs, path)
				}
			}

			if len(foundConfigs) == 0 {
				fmt.Println("No configuration file found to clear.")
				return nil
			}

			// If multiple configs found, list them
			if len(foundConfigs) > 1 {
				fmt.Println("Multiple configuration files found:")
				for _, config := range foundConfigs {
					fmt.Printf("  - %s\n", config)
				}
				if !force {
					fmt.Println("\nUse --force to clear all found configuration files.")
					return nil
				}
			}

			// Clear all found configs
			for _, config := range foundConfigs {
				if err := clearConfigFile(config, force); err != nil {
					return err
				}
			}
			return nil
		}

		// Clear the currently loaded config file
		return clearConfigFile(configFile, force)
	},
}

// clearConfigFile removes a configuration file with optional confirmation
func clearConfigFile(configFile string, force bool) error {
	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Configuration file %s does not exist.\n", configFile)
		return nil
	}

	// Ask for confirmation unless --force is used
	if !force {
		fmt.Printf("Are you sure you want to delete the configuration file: %s? (y/N): ", configFile)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %v", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Configuration clear cancelled.")
			return nil
		}
	}

	// Remove the configuration file
	if err := os.Remove(configFile); err != nil {
		return fmt.Errorf("failed to remove configuration file %s: %v", configFile, err)
	}

	fmt.Printf("Configuration file %s has been removed.\n", configFile)
	fmt.Println("All settings have been reset to defaults.")
	fmt.Println("Environment variables (HYPON_USER, HYPON_PASS) are not affected.")

	return nil
}

// getConfigValueWithSource returns the config value and indicates if it's set
func getConfigValueWithSource(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		if defaultValue != "" {
			return fmt.Sprintf("%s (default)", defaultValue)
		}
		return "<not set>"
	}
	return value
}

// maskPassword masks password values for security
func maskPassword(password string) string {
	if password == "" {
		return "<not set>"
	}
	return "****** (set)"
}

// maskIfSet shows if a value is set without revealing it
func maskIfSet(value string) string {
	if value == "" {
		return "<not set>"
	}
	return "****** (set)"
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configClearCmd)

	// Add flags for clear command
	configClearCmd.Flags().BoolP("force", "f", false, "Force clear without confirmation prompt")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configViewCmd.Flags().BoolP("show-passwords", "p", false, "Show actual password values (insecure)")
}
