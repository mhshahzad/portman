# Contributing to Portman

Thank you for your interest in contributing to Portman!

## Development Setup

1. **Prerequisites**: Go 1.24+
2. **Clone the repository**:
   ```bash
   git clone https://github.com/mhshahzad/portman.git
   cd portman
   ```
3. **Run tests**:
   ```bash
   go test ./...
   ```

## Code Organization

- `/cmd`: CLI command definitions (Cobra).
- `/internal/ports`: Domain models and business logic.
- `/internal/scanner`: OS-level scanning and parsing logic.
- `/internal/output`: Presentation logic.

## Contribution Guidelines

1. **Keep it simple**: We value minimal dependencies and readable code.
2. **Test your changes**: Ensure all logic is covered by unit tests.
3. **Follow Go idioms**: Use standard formatting (`go fmt`) and naming conventions.
4. **Environment-First**: Remember that this tool is primarily designed for Linux/MacOS environments.

## Reporting Issues

Please use the GitHub issue tracker to report bugs or suggest features.
