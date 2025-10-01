# Hypon API CLI

A command-line interface for interacting with the Hypon Cloud API, built with Go using the Cobra CLI framework and Viper for configuration management.

## Features

- **Cobra CLI Framework**: Provides a powerful command-line interface with subcommands, flags, and help text
- **Viper Configuration**: Supports configuration files, environment variables, and command-line flags
- **Modular Architecture**: Clean separation between CLI commands and API client logic
- **Comprehensive Testing**: Unit tests for all major components
- **Authentication**: Built-in login functionality for Hypon Cloud API
- **Extensible**: Easy to add new commands and API operations

## Installation

### Build from Source

```bash
git clone <repository-url>
cd hypon-api
go build -o hypon-api
```

### Run Tests

```bash
go test ./...
```

## Configuration

The CLI supports multiple configuration methods:

### 1. Configuration File

Create a configuration file at `~/.hypon-api.yaml` or specify with `--config` flag:

```yaml
# User configuration
user:
  name: "Your Name Here"

# Hypon API configuration
hypon:
  api_url: "https://api.hypon.cloud"
  # Optional: credentials (can also use environment variables)
  # username: "your-username"
  # password: "your-password"
  oem: ""

# Default settings
defaults:
  page_size: 10
```

### 2. Environment Variables

```bash
export HYPON_USER="your-username"
export HYPON_PASS="your-password"
```

### 3. Command-line Flags

Most commands support flags to override configuration values.

## Usage

### Basic Commands

```bash
# Show help
./hypon-api --help

# Hello world command (demonstrates basic functionality)
./hypon-api hello
./hypon-api hello --name "Your Name"

# API commands
./hypon-api api --help
./hypon-api api admin-info
```

### Configuration Examples

```bash
# Use custom config file
./hypon-api --config /path/to/config.yaml hello

# Hello with configuration from file
./hypon-api hello  # Uses name from config file if set

# Hello with command-line override
./hypon-api hello --name "CLI User"  # Overrides config file setting
```

## Project Structure

```
.
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command and Viper configuration
│   ├── hello.go           # Sample hello command
│   ├── api.go             # API-related commands
│   ├── root_test.go       # Tests for root command
│   └── hello_test.go      # Tests for hello command
├── internal/              # Internal packages
│   └── client/            # API client library
│       ├── client.go      # Hypon API client implementation
│       └── client_test.go # Tests for API client
├── main.go                # Main entry point
├── main_test.go           # Tests for main package
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── .hypon-api.yaml        # Local configuration file
└── .hypon-api.yaml.example # Example configuration file
```

## API Client

The `internal/client` package provides a reusable client for the Hypon Cloud API:

```go
// Create a new client
client, err := client.NewClient("https://api.hypon.cloud")
if err != nil {
    log.Fatal(err)
}

// Authenticate
err = client.Login(username, password, oem)
if err != nil {
    log.Fatal(err)
}

// Make API calls
resp, err := client.GetAdminInfo()
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
```

## Available Commands

### Root Command
- **Usage**: `hypon-api [command]`
- **Flags**: 
  - `--config string`: Config file path (default: `$HOME/.hypon-api.yaml`)

### Hello Command
- **Usage**: `hypon-api hello [flags]`
- **Description**: Sample command demonstrating CLI functionality
- **Flags**:
  - `--name, -n string`: Name to greet (can also be set in config)

### API Commands
- **Usage**: `hypon-api api [command]`
- **Subcommands**:
  - `admin-info`: Get administrator information from Hypon API

## Testing

The project includes comprehensive tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/client
go test ./cmd
```

Test coverage includes:
- API client functionality with mock HTTP servers
- CLI command parsing and execution
- Configuration loading and flag handling
- Error cases and edge conditions

## Development

### Adding New Commands

1. Create a new file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Define the command using Cobra:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Description of new command",
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
    // Add flags if needed
    newCmd.Flags().StringP("flag", "f", "", "Flag description")
}
```

3. Add tests in `cmd/newcommand_test.go`

### Adding New API Methods

1. Add methods to the `Client` struct in `internal/client/client.go`
2. Add corresponding tests in `internal/client/client_test.go`
3. Create CLI commands that use the new methods

## Dependencies

- **[Cobra](https://github.com/spf13/cobra)**: CLI framework
- **[Viper](https://github.com/spf13/viper)**: Configuration management
- **[Testify](https://github.com/stretchr/testify)**: Testing framework

## License

[Add your license information here]

## Contributing

[Add contributing guidelines here]