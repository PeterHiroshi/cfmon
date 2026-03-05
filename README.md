<div align="center">

# 🔨 forge-me

**A lightweight CLI for Cloudflare resource monitoring**

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/PeterHiroshi/forge-me)](https://github.com/PeterHiroshi/forge-me/releases)

</div>

## 📖 Overview

**forge-me** is a fast, intuitive CLI tool for monitoring and managing your Cloudflare resources. Built for developers who prefer the command line over web dashboards, it provides instant access to your Workers and Containers with detailed resource usage metrics.

## 📑 Table of Contents

- [Overview](#-overview)
- [Quick Start](#-quick-start)
- [Demo](#-demo)
- [Why forge-me?](#-why-forge-me)
- [Features](#-features)
- [How It Works](#-how-it-works)
- [Roadmap](#️-roadmap)
- [Installation](#-installation)
- [Usage](#-usage)
- [Examples & Use Cases](#-examples--use-cases)
- [FAQ & Troubleshooting](#-faq--troubleshooting)
- [Configuration](#️-configuration)
- [Development](#️-development)
- [Project Structure](#project-structure)
- [License](#-license)
- [Contributing](#-contributing)
- [Support & Community](#-support--community)

---

## ⚡ Quick Start

```bash
# Install forge-me
curl -sSL https://raw.githubusercontent.com/PeterHiroshi/forge-me/main/scripts/install.sh | bash

# Save your Cloudflare API token
forge-me login YOUR_API_TOKEN

# List your containers
forge-me containers YOUR_ACCOUNT_ID

# List your workers
forge-me workers YOUR_ACCOUNT_ID

# That's it! 🎉
```

---

## 📸 Demo

### Container Listing

```bash
$ forge-me containers abc123def456

ID                    Name               CPU (ms)  Memory (MB)
--------------------  -----------------  --------  -----------
ctr-prod-web-01       production-web     2450      512
ctr-staging-api-02    staging-api        890       256
ctr-dev-worker-03     dev-worker         125       128
```

### Worker Monitoring

```bash
$ forge-me workers abc123def456

ID         Name              CPU (ms)  Requests
---------  ----------------  --------  --------
wkr-001    api-gateway       5230      125400
wkr-002    image-optimizer   3100      89200
wkr-003    auth-service      1850      45600
```

### JSON Output for Automation

```bash
$ forge-me workers abc123def456 --format json | jq '.[0]'

{
  "id": "wkr-001",
  "name": "api-gateway",
  "cpu_ms": 5230,
  "requests": 125400
}
```

> 💡 **Tip**: Add a demo GIF here to showcase the tool in action!
> ```markdown
> ![Demo](https://raw.githubusercontent.com/PeterHiroshi/forge-me/main/docs/demo.gif)
> ```

---

## 🤔 Why forge-me?

The Cloudflare dashboard is powerful, but sometimes you just need quick answers from your terminal:

- **🚀 Speed First**: Get resource metrics in milliseconds, not page loads
- **⌨️ Developer-Friendly**: Designed for CLI workflows and automation
- **🪶 Lightweight**: No browser tabs, no UI bloat—just the data you need
- **🔄 Scriptable**: JSON output mode for easy integration with other tools
- **🔐 Secure**: Token storage keeps your credentials safe locally
- **📊 Focused**: Shows what matters—CPU, memory, requests—without the noise

Perfect for:
- Quick status checks during development
- CI/CD pipeline monitoring
- Resource usage auditing
- Debugging performance issues
- Automating Cloudflare workflows

### Comparison: Dashboard vs forge-me

| Task | Cloudflare Dashboard | forge-me |
|------|---------------------|----------|
| Check worker CPU usage | 1. Open browser<br>2. Log in<br>3. Navigate to Workers<br>4. Click on worker<br>5. View metrics | `forge-me workers ACCOUNT_ID` |
| Export data for analysis | Manual copy-paste or screenshots | `forge-me workers ACCOUNT_ID --format json > data.json` |
| Automate monitoring | Browser automation (complex) | Simple shell script |
| Check multiple accounts | Switch accounts in UI | `--token` flag per account |
| Time to first result | ~10-30 seconds | ~1 second |

---

## ✨ Features

### 📦 **Containers Management**
- List all containers in your Cloudflare account
- View real-time CPU usage (milliseconds) and memory consumption (MB)
- Quick identification by container ID and name
- Perfect for monitoring Container-as-a-Service workloads

### ⚡ **Workers Monitoring**
- Enumerate all deployed Workers scripts
- Track CPU time consumption per worker
- Monitor request counts and traffic patterns
- Identify performance bottlenecks instantly

### 🎨 **Flexible Output Formats**
- **Table Mode** (default): Clean, human-readable tables for terminal viewing
- **JSON Mode**: Machine-readable output for scripting and automation
- Easy integration with tools like `jq`, `grep`, and custom scripts

### 🔐 **Secure Authentication**
- **Login Command**: Save your API token once, use everywhere
- Encrypted local storage in `~/.forge-me/config.yaml`
- Override with `--token` flag for multi-account workflows
- No credentials stored in command history

### 📊 **Status & Health Checks**
- **Ping Command**: Quick connectivity test to Cloudflare API
- Verify authentication and API reachability
- Useful for troubleshooting and CI/CD health checks

### 🐚 **Shell Completion**
- Auto-completion for commands, flags, and arguments
- Supported shells: Bash, Zsh, Fish, PowerShell
- Faster workflows with tab completion

### 🌍 **Cross-Platform Support**
- Linux (x64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (x64, ARM64)
- FreeBSD and NetBSD support

---

## 🔧 How It Works

**forge-me** is a thin, efficient wrapper around the [Cloudflare API v4](https://developers.cloudflare.com/api/). It handles authentication, API requests, and response formatting so you don't have to.

### Architecture Overview

```
┌─────────────┐      HTTPS/REST      ┌──────────────────┐
│             │ ──────────────────▶  │                  │
│  forge-me   │                      │  Cloudflare API  │
│     CLI     │ ◀──────────────────  │       v4         │
│             │      JSON Response   │                  │
└─────────────┘                      └──────────────────┘
      │
      ▼
┌─────────────────────────────────┐
│  Local Config (~/.forge-me/)    │
│  • Encrypted API token          │
│  • User preferences             │
└─────────────────────────────────┘
```

### Key Components

1. **API Client** (`internal/api/`)
   - Manages HTTP connections to Cloudflare
   - Handles authentication headers and error responses
   - Implements rate limiting and retry logic

2. **Configuration Manager** (`internal/config/`)
   - Securely stores and retrieves API tokens
   - Manages user preferences and defaults
   - Supports custom config file locations

3. **Output Formatters** (`internal/output/`)
   - Transforms API responses into readable tables
   - Provides JSON serialization for automation
   - Handles column alignment and formatting

4. **Command Layer** (`cmd/`)
   - Cobra-based CLI interface
   - Parses flags and arguments
   - Orchestrates API calls and output

### API Endpoints Used

| Resource   | Endpoint | Documentation |
|------------|----------|---------------|
| Containers | `/accounts/{account_id}/workers/containers/namespaces` | [Containers API](https://developers.cloudflare.com/api/operations/workers-for-platforms-containers-list) |
| Workers    | `/accounts/{account_id}/workers/scripts` | [Workers API](https://developers.cloudflare.com/api/operations/worker-script-list-workers) |
| Account    | `/accounts/{account_id}` | [Accounts API](https://developers.cloudflare.com/api/operations/accounts-list-accounts) |

---

## 🗺️ Roadmap

We're actively developing new features to make **forge-me** even more powerful:

### 🔜 Coming Soon

- **📺 Watch Mode**: Live dashboard with auto-refreshing metrics
  - Real-time monitoring of resource usage
  - Configurable refresh intervals
  - Interactive TUI with keyboard shortcuts

- **📜 Logs Command**: Tail worker and container logs directly
  - Stream logs from Cloudflare to your terminal
  - Filter by log level, time range, and keywords
  - Export logs to files for analysis

- **🚀 Deploy Command**: Quick worker deployment from CLI
  - Deploy Workers scripts with a single command
  - Support for environment variables and secrets
  - Rollback and version management

### 🔮 Future Ideas

- **👥 Multi-Account Support**: Switch between accounts seamlessly
- **🔍 Sort & Filter Options**: Advanced querying of resources
- **📈 Historical Metrics**: Track usage trends over time
- **🔔 Alerts & Notifications**: Get notified when thresholds are exceeded
- **🎯 Resource Tagging**: Organize and label your resources

Want to contribute? Check out our [Contributing Guide](#contributing) or [open an issue](https://github.com/PeterHiroshi/forge-me/issues) with your ideas!

---

## 📥 Installation

Choose your preferred installation method:

### 🍺 Homebrew (macOS/Linux) — **Recommended**

```bash
brew tap PeterHiroshi/forge-me
brew install forge-me
```

### 🚀 Quick Install Script

**macOS / Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/PeterHiroshi/forge-me/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/PeterHiroshi/forge-me/main/scripts/install.ps1 | iex
```

### 🪣 Package Managers

**Scoop (Windows):**
```powershell
scoop bucket add forge-me https://github.com/PeterHiroshi/scoop-forge-me
scoop install forge-me
```

**APT (Debian/Ubuntu) — Coming Soon:**
```bash
# Install from releases
wget https://github.com/PeterHiroshi/forge-me/releases/latest/download/forge-me_linux_amd64.deb
sudo dpkg -i forge-me_linux_amd64.deb
```

**RPM (Fedora/RHEL) — Coming Soon:**
```bash
# Install from releases
wget https://github.com/PeterHiroshi/forge-me/releases/latest/download/forge-me_linux_amd64.rpm
sudo rpm -i forge-me_linux_amd64.rpm
```

### 📦 Pre-built Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/PeterHiroshi/forge-me/releases):

```bash
# Example for Linux x64
wget https://github.com/PeterHiroshi/forge-me/releases/latest/download/forge-me_linux_amd64.tar.gz
tar -xzf forge-me_linux_amd64.tar.gz
sudo mv forge-me /usr/local/bin/
```

### 🔨 From Source

Requires Go 1.21+:

```bash
git clone https://github.com/PeterHiroshi/forge-me.git
cd forge-me
go build -o forge-me .
sudo mv forge-me /usr/local/bin/
```

### ✅ Verify Installation

```bash
forge-me version
# Output: forge-me version x.x.x
```

---

## 🚀 Usage

### Getting Your API Token

Before using **forge-me**, you'll need a Cloudflare API token:

1. Log in to the [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. Go to **My Profile** → **API Tokens**
3. Click **Create Token**
4. Use the **Read All Resources** template or create a custom token with:
   - **Permissions**: `Account > Workers Scripts > Read`, `Account > Containers > Read`
   - **Account Resources**: Include the accounts you want to monitor

### Authentication

**Recommended:** Save your API token securely for all future commands:

```bash
forge-me login YOUR_API_TOKEN
```

This saves the token to `~/.forge-me/config.yaml`.

Alternatively, provide the token with each command:

```bash
forge-me containers --token YOUR_API_TOKEN ACCOUNT_ID
```

### List Containers

List all containers for an account:

```bash
forge-me containers ACCOUNT_ID
```

Output in JSON format:

```bash
forge-me containers ACCOUNT_ID --format json
```

Example output (table):

```
ID                    Name               CPU (ms)  Memory (MB)
--------------------  -----------------  --------  -----------
container-1           my-container-1     1000      128
container-2           my-container-2     2000      256
```

Example output (JSON):

```json
[
  {
    "id": "container-1",
    "name": "my-container-1",
    "cpu_ms": 1000,
    "memory_mb": 128
  },
  {
    "id": "container-2",
    "name": "my-container-2",
    "cpu_ms": 2000,
    "memory_mb": 256
  }
]
```

### List Workers

List all workers for an account:

```bash
forge-me workers ACCOUNT_ID
```

Output in JSON format:

```bash
forge-me workers ACCOUNT_ID --format json
```

Example output (table):

```
ID         Name           CPU (ms)  Requests
---------  -------------  --------  --------
worker-1   my-worker-1    500       1000
worker-2   my-worker-2    750       2000
```

### Version

Check the installed version:

```bash
forge-me version
```

### Shell Completion

Generate shell completion scripts for your shell:

**Bash:**

```bash
# Load completions in current session
source <(forge-me completion bash)

# Load completions for all sessions
# Linux:
forge-me completion bash > /etc/bash_completion.d/forge-me
# macOS:
forge-me completion bash > $(brew --prefix)/etc/bash_completion.d/forge-me
```

**Zsh:**

```bash
# Enable completion if not already enabled
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Load completions for all sessions
forge-me completion zsh > "${fpath[1]}/_forge-me"
```

**Fish:**

```bash
# Load completions in current session
forge-me completion fish | source

# Load completions for all sessions
forge-me completion fish > ~/.config/fish/completions/forge-me.fish
```

**PowerShell:**

```powershell
# Load completions in current session
forge-me completion powershell | Out-String | Invoke-Expression

# Load completions for all sessions
forge-me completion powershell > forge-me.ps1
# Then source this file from your PowerShell profile
```

### Help

Get help on any command:

```bash
forge-me --help
forge-me containers --help
forge-me workers --help
```

---

## 💡 Examples & Use Cases

### Quick Health Check

```bash
# Test API connectivity
forge-me ping

# List all resources
forge-me containers YOUR_ACCOUNT_ID
forge-me workers YOUR_ACCOUNT_ID
```

### Automation & Scripting

```bash
# Get JSON output for processing
forge-me workers YOUR_ACCOUNT_ID --format json | jq '.[] | select(.cpu_ms > 1000)'

# Monitor specific container
forge-me containers YOUR_ACCOUNT_ID --format json | jq '.[] | select(.name == "my-container")'

# Check if any worker is using high CPU
HIGH_CPU=$(forge-me workers YOUR_ACCOUNT_ID --format json | jq '[.[] | select(.cpu_ms > 5000)] | length')
if [ "$HIGH_CPU" -gt 0 ]; then
  echo "Warning: High CPU usage detected!"
fi
```

### CI/CD Integration

```bash
# In your CI pipeline
export CLOUDFLARE_API_TOKEN="${{ secrets.CLOUDFLARE_TOKEN }}"
forge-me workers "$ACCOUNT_ID" --format json > workers-status.json

# Validate deployment
forge-me workers "$ACCOUNT_ID" --format json | \
  jq -e '.[] | select(.name == "production-worker")' || exit 1
```

### Multi-Account Management

```bash
# Use different tokens per account
forge-me containers ACCOUNT_1 --token "$TOKEN_1"
forge-me containers ACCOUNT_2 --token "$TOKEN_2"

# Or switch config files
forge-me --config ~/.forge-me/account1.yaml containers ACCOUNT_1
forge-me --config ~/.forge-me/account2.yaml containers ACCOUNT_2
```

---

## ❓ FAQ & Troubleshooting

### Common Issues

**Q: I get "unauthorized" errors**
```bash
# Check if your token is valid
forge-me ping

# Verify token has correct permissions in Cloudflare dashboard
# Required: Workers Scripts:Read, Containers:Read
```

**Q: Where is my account ID?**
- Log in to [Cloudflare Dashboard](https://dash.cloudflare.com/)
- Select your account
- Find the Account ID on the right sidebar

**Q: Can I use multiple accounts?**
- Yes! Use the `--token` flag for each account, or maintain separate config files with `--config`

**Q: JSON output is too verbose**
```bash
# Pipe through jq for filtering
forge-me workers ACCOUNT_ID --format json | jq '.[] | {name, cpu_ms}'
```

**Q: How do I uninstall?**
```bash
# Remove the binary
rm $(which forge-me)

# Remove config (optional)
rm -rf ~/.forge-me
```

### Debug Mode

For troubleshooting, you can enable verbose output:

```bash
# Set log level
export LOG_LEVEL=debug
forge-me containers ACCOUNT_ID
```

---

## ⚙️ Configuration

### Config File Location

Configuration is stored in `~/.forge-me/config.yaml` by default.

**Custom config location:**

```bash
forge-me --config /path/to/config.yaml containers ACCOUNT_ID
```

### Config File Format

```yaml
# Your Cloudflare API token
token: your-cloudflare-api-token

# Optional: Default output format
# format: json

# Optional: Default account ID
# account_id: your-account-id
```

### Environment Variables

You can also use environment variables (they take precedence over config file):

```bash
export CLOUDFLARE_API_TOKEN="your-token"
export CLOUDFLARE_ACCOUNT_ID="your-account-id"

forge-me containers  # Uses env vars
```

### Global Flags

All commands support these global flags:

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format (`table` or `json`) | `table` |
| `--token` | Cloudflare API token (overrides config) | - |
| `--config` | Custom config file path | `~/.forge-me/config.yaml` |
| `--help` | Show help for command | - |
| `--version` | Show version information | - |

---

## 🛠️ Development

### Prerequisites

- **Go 1.21 or later** ([Download](https://go.dev/dl/))
- **Make** (optional, for convenience commands)
- **GoReleaser** (optional, for releases)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/PeterHiroshi/forge-me.git
cd forge-me

# Install dependencies
go mod download

# Build the binary
go build -o forge-me .

# Run locally
./forge-me --help
```

### Development Workflow

```bash
# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests with verbose output
go test ./... -v

# Run specific test
go test -run TestContainersList ./internal/api

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Build for your platform
go build -o forge-me .
```

### Cross-Platform Builds

```bash
# Build for all platforms with GoReleaser
goreleaser build --snapshot --clean

# Manual cross-compilation examples
GOOS=linux GOARCH=amd64 go build -o forge-me-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o forge-me-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o forge-me-windows-amd64.exe .
```

### Architecture Deep Dive

**forge-me** follows a clean, modular architecture:

```
┌─────────────────────────────────────────────┐
│              CLI Layer (cmd/)               │
│  • Command definitions (Cobra)              │
│  • Flag parsing and validation              │
│  • User input handling                      │
└────────────┬────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────┐
│          Business Logic Layer               │
│  • Orchestrates API calls                   │
│  • Handles authentication flow              │
│  • Coordinates formatting                   │
└────────────┬────────────────────────────────┘
             │
      ┌──────┴──────┐
      ▼             ▼
┌───────────┐  ┌───────────┐
│ API Layer │  │ Config    │
│ (internal │  │ (internal │
│ /api/)    │  │ /config/) │
└───────────┘  └───────────┘
      │
      ▼
┌────────────────┐
│ Output Layer   │
│ (internal/     │
│ output/)       │
└────────────────┘
```

### Key Design Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Testability**: Business logic is decoupled from CLI and external dependencies
3. **Error Handling**: Errors bubble up with context, not panic
4. **Configuration**: 12-factor app principles (env vars, config files, flags)
5. **Performance**: Minimal allocations, concurrent API calls where possible

## Project Structure

```
forge-me/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command with global flags
│   ├── login.go           # Login command
│   ├── containers.go      # Containers command
│   ├── workers.go         # Workers command
│   └── version.go         # Version command
├── internal/
│   ├── api/               # Cloudflare API client
│   │   ├── client.go      # HTTP client
│   │   ├── containers.go  # Container endpoints
│   │   └── workers.go     # Worker endpoints
│   ├── config/            # Configuration management
│   │   └── config.go      # Load/save config
│   └── output/            # Output formatting
│       └── formatter.go   # Table and JSON formatters
├── scripts/               # Install scripts
│   ├── install.sh         # macOS/Linux installer
│   └── install.ps1        # Windows installer
├── main.go                # Entry point
├── go.mod                 # Go dependencies
├── .goreleaser.yaml       # GoReleaser config
└── README.md              # This file
```

---

## 📜 License

**forge-me** is open source software licensed under the [MIT License](LICENSE).

```
MIT License

Copyright (c) 2026 PeterHiroshi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software...
```

See the [LICENSE](LICENSE) file for the full license text.

---

## 🤝 Contributing

We welcome contributions from the community! Whether you're fixing bugs, adding features, or improving documentation, your help is appreciated.

### How to Contribute

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
   ```bash
   git clone https://github.com/YOUR_USERNAME/forge-me.git
   cd forge-me
   ```
3. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
4. **Make your changes** and commit them
   ```bash
   git commit -m "Add amazing feature"
   ```
5. **Push to your fork**
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Open a Pull Request** on GitHub

### Contribution Guidelines

- Write clear, descriptive commit messages
- Add tests for new features
- Update documentation as needed
- Follow existing code style (run `go fmt`)
- Ensure all tests pass (`go test ./...`)
- Keep PRs focused on a single feature or fix

### Areas We Need Help

- 📝 **Documentation**: Tutorials, examples, API docs
- 🐛 **Bug Reports**: Found a bug? Open an issue!
- ✨ **Feature Requests**: Have an idea? We'd love to hear it
- 🧪 **Testing**: Improve test coverage
- 🌍 **Localization**: Help translate error messages

---

## 💬 Support & Community

### Getting Help

- 📖 **Documentation**: You're reading it! Check the sections above
- 🐛 **Bug Reports**: [Open an issue](https://github.com/PeterHiroshi/forge-me/issues/new?template=bug_report.md)
- 💡 **Feature Requests**: [Request a feature](https://github.com/PeterHiroshi/forge-me/issues/new?template=feature_request.md)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/PeterHiroshi/forge-me/discussions)

### Reporting Issues

When reporting bugs, please include:

- **forge-me version** (`forge-me version`)
- **Operating system** and architecture
- **Go version** (if building from source)
- **Steps to reproduce** the issue
- **Expected vs actual behavior**
- **Relevant logs** or error messages

### Security Vulnerabilities

If you discover a security vulnerability, please **DO NOT** open a public issue. Instead, email the maintainer directly at [security@example.com](mailto:security@example.com).

---

## 🙏 Acknowledgments

Special thanks to:

- **[Cloudflare](https://cloudflare.com)** for their excellent API and developer tools
- **[Cobra](https://github.com/spf13/cobra)** for the CLI framework
- **[Viper](https://github.com/spf13/viper)** for configuration management
- **[GoReleaser](https://goreleaser.com)** for seamless multi-platform releases
- All our [contributors](https://github.com/PeterHiroshi/forge-me/graphs/contributors) who help improve forge-me

### Similar Projects

- **[Wrangler](https://github.com/cloudflare/wrangler2)** - Official Cloudflare CLI (focused on deployment)
- **[cf-tool](https://github.com/xalanq/cf-tool)** - Codeforces CLI tool
- **[cloudflare-cli](https://github.com/danielpigott/cloudflare-cli)** - Ruby-based Cloudflare CLI

### Why forge-me is Different

| Feature | forge-me | Wrangler | cloudflare-cli |
|---------|----------|----------|----------------|
| 🎯 Focus | Resource monitoring | Deployment & dev | DNS/Zone management |
| 🪶 Binary Size | ~10MB | ~50MB | Requires Ruby runtime |
| ⚡ Speed | Instant | Fast | Slower (interpreted) |
| 📊 Output Formats | Table + JSON | Text | Text |
| 🔐 Token Storage | ✅ | ❌ | ❌ |
| 🐚 Shell Completion | ✅ | ✅ | ❌ |
| 📦 Standalone Binary | ✅ | ✅ | ❌ (requires Ruby) |

---

## ⭐ Star History

If you find **forge-me** useful, please consider giving it a star on GitHub! It helps others discover the project.

[![Star History Chart](https://api.star-history.com/svg?repos=PeterHiroshi/forge-me&type=Date)](https://star-history.com/#PeterHiroshi/forge-me&Date)

---

<div align="center">

**Built with ❤️ by developers, for developers**

[⭐ Star on GitHub](https://github.com/PeterHiroshi/forge-me) • [🐛 Report Bug](https://github.com/PeterHiroshi/forge-me/issues) • [💡 Request Feature](https://github.com/PeterHiroshi/forge-me/issues)

**Happy forging! 🔨**

</div>