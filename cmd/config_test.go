package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestConfigCmd(t *testing.T) {
	// Test that config command exists and has correct metadata
	if configCmd.Use != "config" {
		t.Errorf("Expected config command Use to be 'config', got '%s'", configCmd.Use)
	}

	if configCmd.Short == "" {
		t.Error("Expected config command to have a Short description")
	}

	if configCmd.Long == "" {
		t.Error("Expected config command to have a Long description")
	}
}

func TestConfigViewCmd(t *testing.T) {
	// Test that config view command exists and has correct metadata
	if configViewCmd.Use != "view" {
		t.Errorf("Expected config view command Use to be 'view', got '%s'", configViewCmd.Use)
	}

	if configViewCmd.Short == "" {
		t.Error("Expected config view command to have a Short description")
	}

	if configViewCmd.Long == "" {
		t.Error("Expected config view command to have a Long description")
	}

	if configViewCmd.Run == nil {
		t.Error("Expected config view command to have a Run function")
	}
}

func TestConfigViewOutput(t *testing.T) {
	// Reset viper state
	viper.Reset()

	// Capture output
	var buf bytes.Buffer

	// Create a new root command for testing to avoid side effects
	testRootCmd := &cobra.Command{Use: "test"}
	testConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management commands",
	}
	testConfigViewCmd := &cobra.Command{
		Use:   "view",
		Short: "View current configuration",
		Run:   configViewCmd.Run, // Use the actual implementation
	}

	testRootCmd.AddCommand(testConfigCmd)
	testConfigCmd.AddCommand(testConfigViewCmd)
	testRootCmd.SetOut(&buf)
	testRootCmd.SetErr(&buf)

	// Execute the command
	testRootCmd.SetArgs([]string{"config", "view"})
	err := testRootCmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	output := buf.String()

	// Verify expected sections are present
	expectedSections := []string{
		"=== Hypon API Configuration ===",
		"Configuration file:",
		"Hypon API Settings:",
		"api_url:",
		"username:",
		"password:",
		"oem:",
		"User Settings:",
		"name:",
		"Environment Variables:",
		"HYPON_USER:",
		"HYPON_PASS:",
		"Config file search paths:",
	}

	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Expected output to contain '%s', but it was missing", section)
		}
	}
}

func TestConfigViewWithConfigValues(t *testing.T) {
	// Reset viper state
	viper.Reset()

	// Set some test configuration values
	viper.Set("hypon.api_url", "https://test.api.com")
	viper.Set("hypon.username", "testuser")
	viper.Set("user.name", "Test User")

	// Capture output
	var buf bytes.Buffer

	// Create a new root command for testing
	testRootCmd := &cobra.Command{Use: "test"}
	testConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management commands",
	}
	testConfigViewCmd := &cobra.Command{
		Use:   "view",
		Short: "View current configuration",
		Run:   configViewCmd.Run,
	}

	testRootCmd.AddCommand(testConfigCmd)
	testConfigCmd.AddCommand(testConfigViewCmd)
	testRootCmd.SetOut(&buf)
	testRootCmd.SetErr(&buf)

	// Execute the command
	testRootCmd.SetArgs([]string{"config", "view"})
	err := testRootCmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	output := buf.String()

	// Verify configured values appear
	if !strings.Contains(output, "https://test.api.com") {
		t.Error("Expected output to contain the configured API URL")
	}

	if !strings.Contains(output, "testuser") {
		t.Error("Expected output to contain the configured username")
	}

	if !strings.Contains(output, "Test User") {
		t.Error("Expected output to contain the configured user name")
	}
}

func TestConfigViewWithEnvironmentVariables(t *testing.T) {
	// Reset viper state
	viper.Reset()

	// Set environment variables
	os.Setenv("HYPON_USER", "envuser")
	os.Setenv("HYPON_PASS", "envpass")
	defer func() {
		os.Unsetenv("HYPON_USER")
		os.Unsetenv("HYPON_PASS")
	}()

	// Capture output
	var buf bytes.Buffer

	// Create a new root command for testing
	testRootCmd := &cobra.Command{Use: "test"}
	testConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management commands",
	}
	testConfigViewCmd := &cobra.Command{
		Use:   "view",
		Short: "View current configuration",
		Run:   configViewCmd.Run,
	}

	testRootCmd.AddCommand(testConfigCmd)
	testConfigCmd.AddCommand(testConfigViewCmd)
	testRootCmd.SetOut(&buf)
	testRootCmd.SetErr(&buf)

	// Execute the command
	testRootCmd.SetArgs([]string{"config", "view"})
	err := testRootCmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	output := buf.String()

	// Verify environment variables are shown as set (masked)
	if !strings.Contains(output, "****** (set)") {
		t.Error("Expected output to show environment variables as set but masked")
	}
}

func TestGetConfigValueWithSource(t *testing.T) {
	// Reset viper state
	viper.Reset()

	tests := []struct {
		name         string
		key          string
		defaultValue string
		configValue  string
		expected     string
	}{
		{
			name:         "config value set",
			key:          "test.key",
			defaultValue: "default",
			configValue:  "configured",
			expected:     "configured",
		},
		{
			name:         "config value not set, has default",
			key:          "test.missing",
			defaultValue: "default",
			configValue:  "",
			expected:     "default (default)",
		},
		{
			name:         "config value not set, no default",
			key:          "test.missing",
			defaultValue: "",
			configValue:  "",
			expected:     "<not set>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			if tt.configValue != "" {
				viper.Set(tt.key, tt.configValue)
			}

			result := getConfigValueWithSource(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestMaskPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{
			name:     "empty password",
			password: "",
			expected: "<not set>",
		},
		{
			name:     "set password",
			password: "secret123",
			expected: "****** (set)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskPassword(tt.password)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestMaskIfSet(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "empty value",
			value:    "",
			expected: "<not set>",
		},
		{
			name:     "set value",
			value:    "somevalue",
			expected: "****** (set)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskIfSet(tt.value)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestConfigClearCmd(t *testing.T) {
	// Test that config clear command exists and has correct metadata
	if configClearCmd.Use != "clear" {
		t.Errorf("Expected config clear command Use to be 'clear', got '%s'", configClearCmd.Use)
	}

	if configClearCmd.Short == "" {
		t.Error("Expected config clear command to have a Short description")
	}

	if configClearCmd.Long == "" {
		t.Error("Expected config clear command to have a Long description")
	}

	if configClearCmd.RunE == nil {
		t.Error("Expected config clear command to have a RunE function")
	}
}

func TestConfigClearCmdFlags(t *testing.T) {
	// Test that the force flag exists
	forceFlag := configClearCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("Expected config clear command to have a 'force' flag")
	}

	if forceFlag.Shorthand != "f" {
		t.Errorf("Expected force flag shorthand to be 'f', got '%s'", forceFlag.Shorthand)
	}
}

func TestClearConfigFileFunction(t *testing.T) {
	// Create a temporary config file for testing
	tmpDir := t.TempDir()
	testConfigFile := tmpDir + "/test-config.yaml"

	// Create test config content
	configContent := `user:
  name: "Test User"
hypon:
  api_url: "https://test.api.com"
`

	err := os.WriteFile(testConfigFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test force clear (no confirmation needed)
	err = clearConfigFile(testConfigFile, true)
	if err != nil {
		t.Errorf("clearConfigFile with force failed: %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(testConfigFile); !os.IsNotExist(err) {
		t.Error("Expected config file to be deleted, but it still exists")
	}
}

func TestClearConfigFileNonExistent(t *testing.T) {
	// Test clearing a non-existent file
	nonExistentFile := "/tmp/non-existent-config.yaml"

	err := clearConfigFile(nonExistentFile, true)
	if err != nil {
		t.Errorf("clearConfigFile should handle non-existent files gracefully, got error: %v", err)
	}
}
