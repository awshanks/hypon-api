package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"awshanks/hypon-api/internal/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication management commands",
	Long: `Commands for managing authentication credentials.
This includes login to save username and password to configuration.`,
}

// authLoginCmd represents the auth login command
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login and save credentials with API authentication",
	Long: `Prompts for username, password, and OEM, then authenticates with the Hypon API.
Successfully retrieved authentication token and credentials are saved to the configuration file.
The stored token can be used for subsequent API calls automatically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performLogin()
	},
}

// performLogin handles the login process
func performLogin() error {
	fmt.Println("=== Hypon API Login ===")
	fmt.Println()

	// Get username
	fmt.Print("Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %v", err)
	}
	username = strings.TrimSpace(username)

	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Get password securely
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %v", err)
	}
	password := string(passwordBytes)
	fmt.Println() // Add newline after password input

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Get OEM (optional)
	fmt.Print("OEM (optional, press Enter to skip): ")
	oem, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read OEM: %v", err)
	}
	oem = strings.TrimSpace(oem)

	fmt.Println()
	fmt.Println("Authenticating with Hypon API...")

	// Get API URL from config or use default
	apiURL := viper.GetString("hypon.api_url")
	if apiURL == "" {
		apiURL = "https://api.hypon.cloud"
	}

	// Create API client and authenticate
	apiClient, err := client.NewClient(apiURL)
	if err != nil {
		return fmt.Errorf("failed to create API client: %v", err)
	}

	// Attempt to login to the API
	if err := apiClient.Login(username, password, oem); err != nil {
		return fmt.Errorf("API authentication failed: %v", err)
	}

	fmt.Println("✓ API authentication successful!")

	// Save credentials and token to config
	if err := saveCredentials(username, password, oem, apiClient.Token); err != nil {
		return fmt.Errorf("failed to save credentials: %v", err)
	}

	fmt.Println()
	fmt.Println("✓ Credentials saved successfully!")
	fmt.Printf("✓ Username: %s\n", username)
	fmt.Println("✓ Password: ****** (saved)")
	if oem != "" {
		fmt.Printf("✓ OEM: %s\n", oem)
	}
	fmt.Println("✓ API Token: ****** (saved)")
	fmt.Println()
	fmt.Println("Your credentials and authentication token have been saved to the configuration file.")
	fmt.Println("You can now use other commands that require authentication.")

	return nil
}

// saveCredentials saves the username, password, OEM, and token to the config file
func saveCredentials(username, password, oem, token string) error {
	// Set the values in viper
	viper.Set("hypon.username", username)
	viper.Set("hypon.password", password)
	if oem != "" {
		viper.Set("hypon.oem", oem)
	}
	viper.Set("hypon.token", token)

	// Determine config file location
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// No config file is loaded, create one
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %v", err)
		}
		configFile = filepath.Join(home, ".hypon-api.yaml")
		viper.SetConfigFile(configFile)
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Write the config file
	if err := viper.WriteConfig(); err != nil {
		// If the file doesn't exist, use WriteConfigAs
		if err := viper.WriteConfigAs(configFile); err != nil {
			return fmt.Errorf("failed to write config file: %v", err)
		}
	}

	fmt.Printf("Configuration saved to: %s\n", configFile)
	return nil
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authLoginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
