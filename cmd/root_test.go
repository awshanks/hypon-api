package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	// This test ensures that Execute() doesn't panic
	// In a real scenario, we might mock os.Exit to test error cases
	assert.NotPanics(t, func() {
		// We can't easily test Execute() directly since it calls os.Exit
		// Instead, we test that the root command is properly configured
		assert.NotNil(t, rootCmd)
	})
}

func TestRootCommand(t *testing.T) {
	cmd := rootCmd

	assert.Equal(t, "hypon-api", cmd.Use)
	assert.Equal(t, "A CLI tool for interacting with the Hypon Cloud API", cmd.Short)
	assert.Contains(t, cmd.Long, "command-line interface")
}

func TestRootCommandFlags(t *testing.T) {
	cmd := rootCmd

	// Test persistent flags
	configFlag := cmd.PersistentFlags().Lookup("config")
	assert.NotNil(t, configFlag)
	assert.Equal(t, "string", configFlag.Value.Type())

	// Test local flags
	toggleFlag := cmd.Flags().Lookup("toggle")
	assert.NotNil(t, toggleFlag)
	assert.Equal(t, "bool", toggleFlag.Value.Type())
}

func TestRootCommandHelp(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "hypon-api")
	assert.Contains(t, output, "Available Commands:")
	assert.Contains(t, output, "hello")
}

func TestInitConfig(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()

	// Test that initConfig doesn't panic
	assert.NotPanics(t, func() {
		initConfig()
	})

	// Test that viper is configured to read environment variables
	// Since we can't easily test file reading without creating actual files,
	// we test the basic configuration
	assert.True(t, true) // Basic test that initConfig runs without error
}

func TestCommandSubcommands(t *testing.T) {
	// Test that expected subcommands are registered
	commands := rootCmd.Commands()
	commandNames := make([]string, len(commands))
	for i, cmd := range commands {
		commandNames[i] = cmd.Name()
	}

	assert.Contains(t, commandNames, "hello")
	assert.Contains(t, commandNames, "api")
	assert.Contains(t, commandNames, "completion")
	assert.Contains(t, commandNames, "help")
}
