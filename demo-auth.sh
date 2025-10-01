#!/bin/bash

# Demo script for the enhanced auth command with API integration
echo "=== Hypon API CLI Enhanced Auth Command Demo ==="
echo ""

echo "1. First, let's see the help for the auth command:"
./hypon-api auth --help
echo ""

echo "2. Now let's see the help for the auth login command:"
./hypon-api auth login --help
echo ""

echo "3. The enhanced auth login command now:"
echo "   - Prompts for username and password"
echo "   - Prompts for OEM (optional)"
echo "   - Authenticates with the Hypon API"
echo "   - Retrieves and stores the authentication token"
echo "   - Saves all credentials to the configuration file"
echo ""

echo "4. After successful login, you can view your saved credentials with:"
echo "   ./hypon-api config view"
echo ""

echo "5. The stored configuration will include:"
echo "   - Username"
echo "   - Password (masked)"
echo "   - OEM (if provided)"
echo "   - API Authentication Token (masked)"
echo ""

echo "6. The token can then be used for subsequent API calls automatically."
echo ""

echo "Demo complete! The enhanced auth command with API integration is ready to use."