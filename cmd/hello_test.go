package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloCommand(t *testing.T) {
	cmd := helloCmd
	
	assert.Equal(t, "hello", cmd.Use)
	assert.Equal(t, "A simple hello world command", cmd.Short)
	assert.Contains(t, cmd.Long, "sample command")
}

func TestHelloCommandFlags(t *testing.T) {
	cmd := helloCmd
	
	nameFlag := cmd.Flags().Lookup("name")
	assert.NotNil(t, nameFlag)
	assert.Equal(t, "string", nameFlag.Value.Type())
	assert.Equal(t, "n", nameFlag.Shorthand)
}

func TestHelloCommandExecution_DefaultName(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()
	
	cmd := helloCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	
	err := cmd.Execute()
	require.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "Hello, World!")
}

func TestHelloCommandExecution_WithNameFlag(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()
	
	cmd := helloCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--name", "TestUser"})
	
	err := cmd.Execute()
	require.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "Hello, TestUser!")
}

func TestHelloCommandExecution_WithViperConfig(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()
	
	// Set config value
	viper.Set("user.name", "ConfigUser")
	
	cmd := helloCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	
	err := cmd.Execute()
	require.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "Hello, ConfigUser!")
}

func TestHelloCommandExecution_FlagOverridesConfig(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()
	
	// Set config value
	viper.Set("user.name", "ConfigUser")
	
	cmd := helloCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--name", "FlagUser"})
	
	err := cmd.Execute()
	require.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "Hello, FlagUser!")
}

func TestHelloCommandExecution_WithAPIConfig(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()
	
	// Set API URL config
	viper.Set("hypon.api_url", "https://test.api.com")
	
	cmd := helloCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	
	err := cmd.Execute()
	require.NoError(t, err)
	
	output := buf.String()
	assert.Contains(t, output, "Hello, World!")
	assert.Contains(t, output, "Hypon API URL configured: https://test.api.com")
}