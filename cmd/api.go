package cmd

import (
	"fmt"
	"os"

	"awshanks/hypon-api/internal/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Interact with the Hypon API",
	Long: `Commands for interacting with the Hypon Cloud API.
This includes authentication and various API operations.`,
}

// adminInfoCmd represents the admin-info command
var adminInfoCmd = &cobra.Command{
	Use:   "admin-info",
	Short: "Get administrator information",
	Long:  `Retrieve administrator information from the Hypon API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get configuration values
		apiURL := viper.GetString("hypon.api_url")
		if apiURL == "" {
			apiURL = "https://api.hypon.cloud"
		}

		username := viper.GetString("hypon.username")
		if username == "" {
			username = os.Getenv("HYPON_USER")
		}

		password := viper.GetString("hypon.password")
		if password == "" {
			password = os.Getenv("HYPON_PASS")
		}

		oem := viper.GetString("hypon.oem")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required. Set them in config file or environment variables HYPON_USER and HYPON_PASS")
		}

		// Create client and authenticate
		client, err := client.NewClient(apiURL)
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}

		if err := client.Login(username, password, oem); err != nil {
			return fmt.Errorf("login failed: %v", err)
		}

		// Get admin info
		resp, err := client.GetAdminInfo()
		if err != nil {
			return fmt.Errorf("failed to get admin info: %v", err)
		}
		defer resp.Body.Close()

		fmt.Printf("Admin info retrieved successfully. Status: %s\n", resp.Status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.AddCommand(adminInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
