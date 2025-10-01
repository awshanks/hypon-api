package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A simple hello world command",
	Long: `A sample command that demonstrates basic CLI functionality.
This command can greet you with a custom message or use configuration
from Viper to personalize the greeting.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = viper.GetString("user.name")
		}
		if name == "" {
			name = "World"
		}

		greeting := fmt.Sprintf("Hello, %s!", name)

		// Check if we have any Hypon API configuration
		apiURL := viper.GetString("hypon.api_url")
		if apiURL != "" {
			greeting += fmt.Sprintf("\nHypon API URL configured: %s", apiURL)
		}

		fmt.Println(greeting)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	helloCmd.Flags().StringP("name", "n", "", "Name to greet (can also be set in config)")
}
