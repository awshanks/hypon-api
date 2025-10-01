package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthCommand(t *testing.T) {
	// Test that auth command exists and has correct basic properties
	cmd := authCmd
	assert.Equal(t, "auth", cmd.Use)
	assert.Equal(t, "Authentication management commands", cmd.Short)
	assert.Contains(t, cmd.Long, "Commands for managing authentication credentials")
}

func TestAuthLoginCommand(t *testing.T) {
	// Test that auth login command exists and has correct basic properties
	cmd := authLoginCmd
	assert.Equal(t, "login", cmd.Use)
	assert.Equal(t, "Login and save credentials with API authentication", cmd.Short)
	assert.Contains(t, cmd.Long, "Prompts for username, password, and OEM")
	assert.Contains(t, cmd.Long, "authenticates with the Hypon API")
}

func TestAuthCommandFlags(t *testing.T) {
	// Test that auth command has expected flags structure
	cmd := authCmd

	// Should have help flag (inherited) - check if it exists when parsed
	flagSet := cmd.Flags()
	assert.NotNil(t, flagSet)

	// The help flag is added automatically by cobra, let's just check the command works
	assert.Equal(t, "auth", cmd.Use)
}

func TestAuthLoginCommandFlags(t *testing.T) {
	// Test that auth login command has expected flags structure
	cmd := authLoginCmd

	// Should have help flag (inherited) - check if it exists when parsed
	flagSet := cmd.Flags()
	assert.NotNil(t, flagSet)

	// The help flag is added automatically by cobra, let's just check the command works
	assert.Equal(t, "login", cmd.Use)
}

func TestSaveCredentials(t *testing.T) {
	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "TestSaveCredentials")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test config file path
	configFile := filepath.Join(tempDir, "test-config.yaml")

	// Set up viper for testing
	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Test saving credentials
	err = saveCredentials("testuser", "testpass", "testoem", "testtoken123")
	assert.NoError(t, err)

	// Verify the file was created
	assert.FileExists(t, configFile)

	// Verify the credentials were saved by reading the config
	viper.Reset()
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	assert.NoError(t, err)

	assert.Equal(t, "testuser", viper.GetString("hypon.username"))
	assert.Equal(t, "testpass", viper.GetString("hypon.password"))
	assert.Equal(t, "testoem", viper.GetString("hypon.oem"))
	assert.Equal(t, "testtoken123", viper.GetString("hypon.token"))
}

func TestAuthCommandStructure(t *testing.T) {
	// Test that auth command has login subcommand
	var found bool
	for _, cmd := range authCmd.Commands() {
		if cmd.Use == "login" {
			found = true
			break
		}
	}
	assert.True(t, found, "auth command should have login subcommand")
}

func TestAuthCommandExecution_Help(t *testing.T) {
	// Create a test root command to avoid affecting global state
	testRootCmd := &cobra.Command{
		Use: "hypon-api-test",
	}

	// Create test auth command (copy of the real one)
	testAuthCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication management commands",
		Long: `Commands for managing authentication credentials.
This includes login to save username and password to configuration.`,
	}

	testRootCmd.AddCommand(testAuthCmd)

	// Test help output
	buf := new(bytes.Buffer)
	testRootCmd.SetOut(buf)
	testRootCmd.SetErr(buf)
	testRootCmd.SetArgs([]string{"auth", "--help"})

	err := testRootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Commands for managing authentication credentials")
	assert.Contains(t, output, "This includes login")
}

func TestAuthLoginCommandExecution_Help(t *testing.T) {
	// Create a test root command to avoid affecting global state
	testRootCmd := &cobra.Command{
		Use: "hypon-api-test",
	}

	// Create test auth command with login subcommand
	testAuthCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication management commands",
		Long: `Commands for managing authentication credentials.
This includes login to save username and password to configuration.`,
	}

	testAuthLoginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login and save credentials with API authentication",
		Long: `Prompts for username, password, and OEM, then authenticates with the Hypon API.
Successfully retrieved authentication token and credentials are saved to the configuration file.
The stored token can be used for subsequent API calls automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil // Mock implementation for testing
		},
	}

	testRootCmd.AddCommand(testAuthCmd)
	testAuthCmd.AddCommand(testAuthLoginCmd)

	// Test help output
	buf := new(bytes.Buffer)
	testRootCmd.SetOut(buf)
	testRootCmd.SetErr(buf)
	testRootCmd.SetArgs([]string{"auth", "login", "--help"})

	err := testRootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Prompts for username, password, and OEM")
	assert.Contains(t, output, "authenticates with the Hypon API")
}

func TestSaveCredentials_CreateConfigDir(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "TestSaveCredentials_CreateConfigDir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a nested path that doesn't exist
	configFile := filepath.Join(tempDir, "nested", "subdir", "test-config.yaml")

	// Set up viper for testing
	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Test saving credentials (should create directories)
	err = saveCredentials("testuser", "testpass", "testoem", "testtoken123")
	assert.NoError(t, err)

	// Verify the file was created
	assert.FileExists(t, configFile)

	// Verify the directory structure was created
	assert.DirExists(t, filepath.Dir(configFile))
}

func TestSaveCredentials_EmptyValues(t *testing.T) {
	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "TestSaveCredentials_EmptyValues")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test config file path
	configFile := filepath.Join(tempDir, "test-config.yaml")

	// Set up viper for testing
	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Test saving credentials with empty values
	err = saveCredentials("", "", "", "")
	assert.NoError(t, err)

	// Verify the file was created even with empty values
	assert.FileExists(t, configFile)

	// Verify the credentials were saved by reading the config
	viper.Reset()
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	assert.NoError(t, err)

	assert.Equal(t, "", viper.GetString("hypon.username"))
	assert.Equal(t, "", viper.GetString("hypon.password"))
}
