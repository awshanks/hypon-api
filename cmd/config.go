package cmd

import (
	"fmt"
	"os"

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configViewCmd.Flags().BoolP("show-passwords", "p", false, "Show actual password values (insecure)")
}