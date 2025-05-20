# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the MCP Server project.

## Adding New Workflows

When adding new workflows, please follow these guidelines:

1. Use appropriate job dependencies to optimize performance
2. Set up proper caching for Go modules and build artifacts
3. Keep workflow files focused on specific tasks
4. Use the standard Go version defined in the Makefile (1.24.x)

## Testing Locally

To test these workflows locally before pushing, you can use [act](https://github.com/nektos/act):

```bash
# Install act
brew install act

# Run a specific workflow
act -j lint
act -j build
act -j test
```