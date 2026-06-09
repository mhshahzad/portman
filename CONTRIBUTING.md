# Contributing to Portman

Thank you for your interest in contributing to Portman!

Be it your first experiment with a CLI tool or your first Go project—you are more than welcome to contribute. Portman is built to be simple, practical, and beginner-friendly for anyone curious about DevOps tooling or Go development.

---

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

---

## Code Organization

* `/cmd`: CLI command definitions (Cobra)
* `/internal/ports`: Domain models and business logic
* `/internal/scanner`: OS-level scanning and parsing logic
* `/internal/output`: Presentation logic

---

## Contribution Guidelines

1. **Keep it simple**
   We value minimal dependencies, clarity, and readable code over complexity.

2. **Test your changes**
   Ensure new logic is covered with unit tests where applicable.

3. **Follow Go idioms**
   Use `go fmt`, standard naming conventions, and idiomatic Go patterns.

4. **Environment-first design**
   This tool is primarily designed for Linux/macOS environments commonly used in DevOps workflows.

---

## Reporting Issues

Please use the GitHub issue tracker to report bugs, suggest improvements, or request features.
