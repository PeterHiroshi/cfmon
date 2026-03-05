# Contributing to cfmon

Thank you for your interest in contributing to cfmon! We welcome contributions from the community and are grateful for your support.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Security](#security)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Accept feedback gracefully
- Put the project's best interests first

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/cfmon.git
   cd cfmon
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/PeterHiroshi/cfmon.git
   ```
4. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Make (for using the Makefile)
- Git

### Installation

1. **Install dependencies**:
   ```bash
   make deps
   ```

2. **Build the project**:
   ```bash
   make build
   ```

3. **Run tests**:
   ```bash
   make test
   ```

### Optional Tools

- **golangci-lint** for linting:
  ```bash
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
  ```

- **air** for hot reload during development:
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

- **goreleaser** for building releases:
  ```bash
  go install github.com/goreleaser/goreleaser@latest
  ```

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug fixes**: Fix issues reported in GitHub Issues
- **Features**: Implement new features or enhance existing ones
- **Documentation**: Improve README, add examples, fix typos
- **Tests**: Add missing tests, improve test coverage
- **Performance**: Optimize code for better performance
- **Refactoring**: Improve code quality and maintainability

### Finding Something to Work On

- Check [GitHub Issues](https://github.com/PeterHiroshi/cfmon/issues)
- Look for issues labeled `good first issue` or `help wanted`
- Propose new features by opening an issue first

## Development Workflow

### 1. Keep Your Fork Updated

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

### 2. Make Your Changes

- Write clean, readable code
- Follow the existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Run Tests and Checks

```bash
# Format your code
make fmt

# Run linters
make lint

# Run tests
make test

# Run all checks
make check
```

### 4. Commit Your Changes

Follow our [commit message conventions](#commit-messages).

### 5. Push and Create a Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Testing

### Unit Tests

- Write tests for all new functions and methods
- Place tests in `*_test.go` files in the same package
- Use table-driven tests where appropriate
- Aim for >90% coverage on new code

Example:

```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "TEST", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewFeature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFeature() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("NewFeature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

- Place integration tests in `test/integration/`
- Use build tags to separate from unit tests
- Test complete workflows and CLI commands

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run integration tests
make integration-test

# Run specific test
go test -v -run TestSpecificFunction ./internal/...
```

## Code Style

### Go Code

- Follow standard Go conventions
- Use `gofmt` for formatting (automatic with `make fmt`)
- Use meaningful variable and function names
- Keep functions small and focused
- Add comments for exported functions and types
- Handle errors explicitly

### Documentation

- Use clear, concise language
- Include code examples where helpful
- Keep README up to date
- Document breaking changes

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("failed to process %s: %w", filename, err)
}

// Bad
if err != nil {
    return err
}
```

## Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code refactoring
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks
- **ci**: CI/CD changes

### Examples

```
feat(containers): add filtering support for container list

Add --filter flag to containers list command to allow filtering
by container name using substring matching.

Closes #123
```

```
fix(auth): handle expired tokens gracefully

Check token expiration before making API calls and provide
helpful error message suggesting to run 'cfmon login' again.
```

## Pull Request Process

### Before Submitting

1. **Update from upstream main**
2. **Run all tests** and ensure they pass
3. **Update documentation** if needed
4. **Add tests** for new functionality
5. **Check code style** with `make lint`

### PR Guidelines

- **Title**: Use a clear, descriptive title
- **Description**: Explain what changes you made and why
- **Testing**: Describe how you tested your changes
- **Screenshots**: Include screenshots for UI changes
- **Breaking changes**: Clearly mark any breaking changes

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] No breaking changes (or documented)
```

### Review Process

- PRs require at least one review
- Address all feedback constructively
- Make requested changes promptly
- Squash commits if requested

## Issue Guidelines

### Bug Reports

Include:
- Clear, descriptive title
- Steps to reproduce
- Expected behavior
- Actual behavior
- System information (OS, Go version)
- cfmon version
- Error messages or logs

### Feature Requests

Include:
- Clear use case
- Proposed solution
- Alternative solutions considered
- Impact on existing functionality

### Issue Template

```markdown
## Description
Clear description of the issue/feature

## Steps to Reproduce (for bugs)
1. Run command '...'
2. See error

## Expected Behavior
What should happen

## Actual Behavior (for bugs)
What actually happens

## Environment
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.21.5]
- cfmon version: [e.g., v0.1.0]
```

## Security

### Reporting Security Issues

**DO NOT** create public issues for security vulnerabilities.

Instead, please email security concerns to the maintainers directly. Include:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will respond within 48 hours and work on a fix.

## Recognition

Contributors will be recognized in:
- The project README
- Release notes
- GitHub contributors page

## Questions?

If you have questions about contributing:

1. Check existing documentation
2. Search closed issues
3. Open a new issue with the `question` label
4. Join our discussions on GitHub Discussions

## License

By contributing to cfmon, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to cfmon! Your efforts help make this project better for everyone. 🎉