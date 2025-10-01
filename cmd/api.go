package cmd

import (
	"encoding/json"
	"fmt"
	"io"
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

		token := viper.GetString("hypon.token")

		// Create client
		client, err := client.NewClient(apiURL)
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}

		// If we have a stored token, use it directly
		if token != "" {
			client.SetToken(token)
		} else {
			// Fall back to username/password authentication
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
				return fmt.Errorf("no stored token found and username/password are required. Run 'hypon-api auth login' first or set HYPON_USER and HYPON_PASS environment variables")
			}

			if err := client.Login(username, password, oem); err != nil {
				return fmt.Errorf("login failed: %v", err)
			}
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

// plantsCmd represents the plants command
var plantsCmd = &cobra.Command{
	Use:   "plants",
	Short: "Get solar plants information",
	Long:  `Retrieve the list of solar plants from the Hypon API using the menu action endpoint.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get configuration values
		apiURL := viper.GetString("hypon.api_url")
		if apiURL == "" {
			apiURL = "https://api.hypon.cloud"
		}

		token := viper.GetString("hypon.token")

		// Create client
		client, err := client.NewClient(apiURL)
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}

		// If we have a stored token, use it directly
		if token != "" {
			client.SetToken(token)
		} else {
			// Fall back to username/password authentication
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
				return fmt.Errorf("no stored token found and username/password are required. Run 'hypon-api auth login' first or set HYPON_USER and HYPON_PASS environment variables")
			}

			if err := client.Login(username, password, oem); err != nil {
				return fmt.Errorf("login failed: %v", err)
			}
		}

		// Get plants list
		resp, err := client.GetPlantListMenu()
		if err != nil {
			return fmt.Errorf("failed to get plants list: %v", err)
		}
		defer resp.Body.Close()

		// Read and display the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		// Try to pretty-print JSON if possible
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
			if err == nil {
				fmt.Printf("Plants list retrieved successfully. Status: %s\n\n", resp.Status)
				fmt.Println(string(prettyJSON))
				return nil
			}
		}

		// Fallback: display raw response
		fmt.Printf("Plants list retrieved successfully. Status: %s\n\n", resp.Status)
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.AddCommand(adminInfoCmd)
	apiCmd.AddCommand(plantsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
